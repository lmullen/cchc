// Package messages sends jobs to a message queue. DEPRECATED.
//
// It provides a Repository interface for generalized interactions storing and
// retrieving messages from a message queue, as well as a concrete type that
// implements that interface for RabbitMQ.
package messages

import (
	"context"

	"github.com/streadway/amqp"
)

// Repository is an interface that sends messages from a particular
// queue in a message broker.
type Repository interface {
	Send(ctx context.Context, body interface{}) error
	Consume() <-chan amqp.Delivery
	Close() error
	QueueName() string
}
