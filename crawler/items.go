package main

import (
	"database/sql"
	"fmt"
	"time"
)

var remove = []string{"more_like_this", "related_items", "cite_this"}

// Item is a representation of an item in the LOC collection returned from the API.
type Item struct {
	ID              string
	URL             string
	Title           string
	Date            time.Time
	Subjects        []string
	Fulltext        string
	FulltextService string
	FulltextFile    string
	Timestamp       int64
	API             []byte // The entire API response stored as JSONB
}

// Exists checks whether an item has been saved to the database or not.
func (i Item) Exists(id string) (bool, error) {
	var exists bool
	err := app.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM items WHERE id=$1)", id).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, fmt.Errorf("Error checking if item %s exists: %w", id, err)
	}
	return exists, nil
}

// ItemResponse represents an item-level object returned from the API. Many more fields
// are returned and will be stored in the database as a JSONB field, but these
// are the ones that will be serialized to regular database fields.
//
// TODO: For now we are keeping track of only information we are sure we are
// going to use. Some of the fields are commented out because they may prove
// useful in the future.
type ItemResponse struct {
	ItemDetails struct {
		ID       string    `json:"id"`
		URL      string    `json:"url"`
		Date     time.Time `json:"date"`
		Subjects []string  `json:"subject_headings"`
		Title    string    `json:"title"`
		// Language     []string  `json:"language"`
		// OnlineFormat []string  `json:"online_format"`
		// Version      int64     `json:"_version_"`
		// HasSegments  bool      `json:"hassegments"`
		// Timestamp time.Time `json:"timestamp"`
	} `json:"item"`
	Resources []struct {
		Files [][]struct {
			Fulltext        string `json:"fulltext,omitempty"`
			FulltextService string `json:"fulltext_service,omitempty"`
		} `json:"files"`
		FulltextFile string `json:"fulltext_file"`
	} `json:"resources"`
	Timestamp int64 `json:"timestamp"`
}
