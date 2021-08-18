package main

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/streadway/amqp"

	log "github.com/sirupsen/logrus"
)

// StartProcessingItems begins reading items from the message queue in order to
// process each of them.
func startProcessingItems(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range app.ItemMetadataQ.Consumer {
		select {
		case <-ctx.Done():
			return
		default:
			// Give each item its own goroutine
			wg.Add(1)
			go processItemMetadata(ctx, wg, msg)
		}
	}
}

// ProcessItemMetadata reads an item from the queue, fetches its metadata, and
// saves it to the database.
func processItemMetadata(ctx context.Context, wg *sync.WaitGroup, msg amqp.Delivery) {
	defer wg.Done()
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
		msg.Ack(false)
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
