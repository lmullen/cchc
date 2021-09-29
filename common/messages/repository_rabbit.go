package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// Repo is a data store using RabbitMQ via the amqp interface.
type Repo struct {
	Channel  *amqp.Channel
	Queue    *amqp.Queue
	Consumer <-chan amqp.Delivery
}

// NewMessageRepo returns a message repo using RabbitMQ via the amqp interface.
func NewMessageRepo(channel *amqp.Channel, queue *amqp.Queue, consumer <-chan amqp.Delivery) *Repo {
	return &Repo{
		Channel:  channel,
		Queue:    queue,
		Consumer: consumer,
	}
}

// Send publishes a message to the message queue
func (r *Repo) Send(ctx context.Context, text *FullTextPredict) error {
	json, err := json.Marshal(text)
	if err != nil {
		return fmt.Errorf("Error marshalling full text to JSON: %w", err)
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
