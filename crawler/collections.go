package main

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

// Collection describes a single collection as returned by the search API.
type Collection struct {
	Contributor      []string  `json:"contributor"`
	Count            int       `json:"count"`
	Description      []string  `json:"description"`
	Digitized        bool      `json:"digitized"`
	ExtractTimestamp time.Time `json:"extract_timestamp"`
	Group            []string  `json:"group"`
	Hassegments      bool      `json:"hassegments"`
	ID               string    `json:"id"`
	ImageURL         []string  `json:"image_url"`
	Index            int       `json:"index"`
	Item             struct {
		AccessAdvisory []string `json:"access_advisory"`
		Contributors   []string `json:"contributors"`
		Date           string   `json:"date"`
		Format         []string `json:"format"`
		Language       []string `json:"language"`
		Location       []string `json:"location"`
		Medium         []string `json:"medium"`
		Notes          []string `json:"notes"`
		Repository     []string `json:"repository"`
		Subjects       []string `json:"subjects"`
		Summary        []string `json:"summary"`
		Title          string   `json:"title"`
	} `json:"item"`
	ItemsURL             string        `json:"items"`
	Language             []string      `json:"language"`
	Location             []string      `json:"location"`
	Number               []string      `json:"number"`
	NumberLccn           []string      `json:"number_lccn"`
	NumberSourceModified []string      `json:"number_source_modified"`
	OriginalFormat       []string      `json:"original_format"`
	OtherTitle           []interface{} `json:"other_title"`
	Partof               []string      `json:"partof"`
	Resources            []interface{} `json:"resources"`
	ShelfID              string        `json:"shelf_id"`
	Subject              []string      `json:"subject"`
	SubjectTopic         []string      `json:"subject_topic"`
	Timestamp            time.Time     `json:"timestamp"`
	Title                string        `json:"title"`
	URL                  string        `json:"url"`
	Dates                []time.Time   `json:"dates,omitempty"`
}

// String prints the title of the digital collection
func (c Collection) String() string {
	return c.Title
}

// Save serializes collection metadata to the database.
func (c Collection) Save() error {
	query := `
	INSERT INTO collections(id, title, description, count, url, items_url, subjects, subjects2, topics, api) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	ON CONFLICT DO NOTHING;
	`

	// Convert the rest of the data back to JSON to stuff into a DB column
	api, err := json.Marshal(c)
	if err != nil {
		log.Debug("Error marshalling JSON to store in collections table", err)
	}

	// TODO: Consider using a single prepared query
	stmt, err := app.DB.Prepare(query)
	if err != nil {
		return err
	}

	// Avoid panicking if the collection does not have a description
	description := ""
	if len(c.Description) > 0 {
		description = c.Description[0]
	}

	_, err = stmt.Exec(c.ID, c.Title, description, c.Count, c.URL, c.ItemsURL, c.Item.Subjects, c.Subject, c.SubjectTopic, api)
	if err != nil {
		return err
	}

	return nil
}
