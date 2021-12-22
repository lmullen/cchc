package messages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/streadway/amqp"
)

// RabbitMQ is a data store using RabbitMQ via its amqp interface.
type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      *amqp.Queue
	Consumer   <-chan amqp.Delivery
}

// connect creates a connection to a RabbitMQ message broker. It will retry the connection several times.
func connect(ctx context.Context, connstr string) (*amqp.Connection, error) {
	var conn *amqp.Connection

	connectWithRetry := func() error {
		select {
		case <-ctx.Done():
			return backoff.Permanent(errors.New("Cancelled attempt to connect to RabbitMQ"))
		default:
			c, err := amqp.Dial(connstr)
			if err != nil {
				return err
			}
			conn = c
		}
		return nil
	}

	policy := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10)
	err := backoff.Retry(connectWithRetry, policy)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to message broker: %w", err)
	}

	return conn, nil
}

// NewRabbitMQ returns a message repo using RabbitMQ via the amqp interface.
// It will create (or connect to) a queue for sending a particular type of message
// (determined by the type of messages sent or receive from the queue) for a
// particular purpose (determined by the name of the queue).
func NewRabbitMQ(ctx context.Context, connstr string, queue string, qos int) (*RabbitMQ, error) {
	repo := RabbitMQ{}
	conn, err := connect(ctx, connstr)
	if err != nil {
		return nil, err
	}
	repo.Connection = conn

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	repo.Channel = ch

	err = ch.Qos(qos, 0, true)
	if err != nil {
		return nil, err
	}

	dle, dlq, dlk := "dead-letter-exchange", "dead-letter-queue", "dead-letter-key"
	err = ch.ExchangeDeclare(dle, "fanout", true, false, false, false, amqp.Table{})
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(dlq, true, false, false, false, amqp.Table{})
	if err != nil {
		return nil, err
	}
	err = ch.QueueBind(dlq, dlk, dle, false, amqp.Table{})
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(queue, true, false, false, false,
		amqp.Table{
			"x-queue-mode":           "lazy",
			"x-dead-letter-exchange": dle,
		})
	if err != nil {
		return nil, err
	}

	repo.Queue = &q

	consumer, err := ch.Consume(q.Name, queue+"-consumer",
		false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	repo.Consumer = consumer

	return &repo, nil
}

// Send publishes a message to the message queue. Body is any kind of object that
// will be marshalled to JSON and included as the body of the message.
func (r *RabbitMQ) Send(ctx context.Context, body interface{}) error {
	json, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("Error marshalling full text message to JSON: %w", err)
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/json",
		Timestamp:    time.Now(),
		Body:         json,
	}

	err = r.Channel.Publish("", r.Queue.Name, false, false, msg)
	if err != nil {
		return fmt.Errorf("Failed to publish full text message: %w", err)
	}

	return nil
}

// Consume returns the channel that provides deliveries from the message queue.
func (r *RabbitMQ) Consume() <-chan amqp.Delivery {
	return r.Consumer
}

// Close shutdowns the connection to the message broker and associated resources.
func (r *RabbitMQ) Close() error {
	err := r.Connection.Close()
	return err
}

func (r *RabbitMQ) QueueName() string {
	return r.Queue.Name
}
