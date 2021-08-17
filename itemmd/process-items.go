package main

import (
	"encoding/json"

	"github.com/streadway/amqp"

	log "github.com/sirupsen/logrus"
)

// StartProcessingItems begins reading items from the message queue in order to
// process each of them.
func startProcessingItems() {
	for msg := range app.ItemMetadataQ.Consumer {
		// Give each item its own goroutine
		go processItemMetadata(msg)
	}
}

// ProcessItemMetadata reads an item from the queue, fetches its metadata, and
// saves it to the database.
func processItemMetadata(msg amqp.Delivery) {
	var item Item
	err := json.Unmarshal(msg.Body, &item)
	if err != nil {
		msg.Reject(false)
		log.WithError(err).WithField("msg", msg).Error("Failed to read body of message from queue")
		return
	}
	// Check if fetched already
	fetched, _ := item.Fetched()
	if fetched {
		log.WithField("id", item.ID).Debug("Skipping item which was queued but already fetched")
		return
	}
	err = item.Fetch()
	if err != nil {
		msg.Reject(false)
		log.WithError(err).WithField("url", item.URL).WithField("id", item.ID).Error("Error fetching item")
		return
	}
	err = item.Save()
	if err != nil {
		msg.Reject(false)
		log.WithError(err).WithField("id", item.ID).Error("Error saving item to database")
		return
	}
	msg.Ack(false)
	return
}
