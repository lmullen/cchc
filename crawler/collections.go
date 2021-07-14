package main

import (
	"encoding/json"
	"time"
)

// Collection describes a single collection
type Collection struct {
	// AccessRestricted bool          `json:"access_restricted"`
	// Aka              []string      `json:"aka"`
	// Campaigns        []interface{} `json:"campaigns"`
	Contributor      []string  `json:"contributor"`
	Count            int       `json:"count"`
	Description      []string  `json:"description"`
	Digitized        bool      `json:"digitized"`
	ExtractTimestamp time.Time `json:"extract_timestamp"`
	Group            []string  `json:"group"`
	// Hassegments      bool          `json:"hassegments"`
	ID string `json:"id"`
	// ImageURL         []string      `json:"image_url"`
	// Index            int           `json:"index"`
	Item struct {
		AccessAdvisory []string    `json:"access_advisory"`
		Contributors   []string    `json:"contributors"`
		Date           string      `json:"date"`
		Format         interface{} `json:"format"`
		Language       []string    `json:"language"`
		Location       []string    `json:"location"`
		Medium         interface{} `json:"medium"`
		Notes          []string    `json:"notes"`
		Repository     []string    `json:"repository"`
		Subjects       []string    `json:"subjects"`
		Summary        []string    `json:"summary"`
		Title          string      `json:"title"`
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
func (cm Collection) String() string {
	return cm.Title
}

// Save serializes collection metadata to the database.
func (cm Collection) Save() error {
	query := `
	INSERT INTO collections(id, url, items_url, count, title, subjects, api) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT DO NOTHING;
	`

	// Convert the rest of the data back to JSON to stuff into a DB column
	//
	// TODO Perhaps this step can be avoided by keeping the unparsed JSON in the struct
	api, _ := json.Marshal(cm)

	stmt, err := app.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(cm.ID, cm.URL, cm.ItemsURL, cm.Count, cm.Title, cm.Item.Subjects, api)
	if err != nil {
		return err
	}

	return nil

}
