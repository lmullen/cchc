package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// ItemResult is an item returned from a LOC.gov collection results page. There
// are many more fields that are returned in a collections result page, but we
// are going to get that data directly from the item page instead.
type ItemResult struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	CollectionID string // This is the foreign key to the collection, not from API
}

// Use the title as a string representation of an item.
func (item ItemResult) String() string {
	return item.ID
}

// Save serializes an item to the database.
func (item ItemResult) Save() error {
	itemQuery := `
	INSERT INTO items(id, url) 
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING;
	`

	relationQuery := `
	INSERT INTO items_in_collections(item_id, collection_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING;
	`
	itemStmt, err := app.DB.Prepare(itemQuery)
	if err != nil {
		return fmt.Errorf("Error preparing item save query: %w", err)
	}

	relationStmt, err := app.DB.Prepare(relationQuery)
	if err != nil {
		return fmt.Errorf("Error preparing item/collection query: %w", err)
	}

	// Use a transaction since we are writing to two tables
	tx, err := app.DB.Begin()
	if err != nil {
		return fmt.Errorf("Error creating transaction in database: %w", err)
	}

	_, err = tx.Stmt(itemStmt).Exec(item.ID, item.URL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error saving item to database: %w", err)
	}

	_, err = tx.Stmt(relationStmt).Exec(item.ID, item.CollectionID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error saving item/collection relation to database: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Error committing item to database: %w", err)
	}

	return nil

}

// Fetched checks whether an item's metadata has already been fetched
// from the API. This will also return false if the item has not been saved at
// all.
func (item ItemResult) Fetched() (bool, error) {
	var fetched bool
	query := `SELECT EXISTS (SELECT 1 FROM items WHERE id=$1 AND api IS NOT NULL)`
	err := app.DB.QueryRow(query, item.ID).Scan(&fetched)
	if err != nil && err != sql.ErrNoRows {
		return fetched, fmt.Errorf("Error checking if item %s has been fetched: %w", item.ID, err)
	}
	return fetched, nil
}

// EnqueueMetadata sends a message to the message queue to so that the item's
// metadata will put into a queue to be fetched.
func (item ItemResult) EnqueueMetadata() error {

	json, err := json.Marshal(item.Msg())
	if err != nil {
		return fmt.Errorf("Error marshalling item to JSON: %w", err)
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/json",
		Timestamp:    time.Now(),
		Body:         json,
	}

	err = app.ItemMetadataQ.Channel.Publish("", app.ItemMetadataQ.Queue.Name,
		false, false, msg)
	if err != nil {
		return fmt.Errorf("Failed to publish item metadata message: %w", err)
	}

	return nil
}

// ItemMetadataMsg is a minimal representation of an item for sending to the
// message broker.
type ItemMetadataMsg struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// Msg returns a minimal representation for sending to the message broker.
func (item ItemResult) Msg() ItemMetadataMsg {
	return ItemMetadataMsg{
		ID:  item.ID,
		URL: item.URL,
	}
}
