package messages

import (
	"context"

	"github.com/streadway/amqp"
)

// Repository is an interface that sends messages from a particular
// queue in a message broker.
type Repository interface {
	Send(ctx context.Context, text *FullTextPredict) error
	Consume() <-chan amqp.Delivery
}
