package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// StartProcessingItems begins reading items from the message queue in order to
// process each of them.
func StartProcessingItems() {
	for msg := range app.ItemMetadataQ.Consumer {
		// Give each item its own goroutine
		go func(msg amqp.Delivery) {
			err := ProcessItemMetadata(msg)
			if err != nil {
				log.Error("Error processing item from queue: ", err)
			}
		}(msg)
	}
}
