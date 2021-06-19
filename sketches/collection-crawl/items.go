package main

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"
)

// ItemResult is an item returned from a LOC.gov collection results page
type ItemResult struct {
	// AccessRestricted bool          `json:"access_restricted"`
	// Aka              []string      `json:"aka"`
	// Campaigns        []interface{} `json:"campaigns"`
	Contributor      []string    `json:"contributor"`
	Date             string      `json:"date"`
	Dates            []time.Time `json:"dates"`
	Description      []string    `json:"description"`
	Digitized        bool        `json:"digitized"`
	ExtractTimestamp time.Time   `json:"extract_timestamp"`
	Group            []string    `json:"group"`
	// Hassegments      bool      `json:"hassegments"`
	ID string `json:"id"`
	// ImageURL         []string  `json:"image_url"`
	// Index            int       `json:"index"`
	Item struct {
		CallNumber       []string `json:"call_number"`
		Contributors     []string `json:"contributors"`
		CreatedPublished []string `json:"created_published"`
		Date             string   `json:"date"`
		Format           []string `json:"format"`
		Language         []string `json:"language"`
		Medium           []string `json:"medium"`
		Notes            []string `json:"notes"`
		OtherTitle       []string `json:"other_title"`
		Subjects         []string `json:"subjects"`
		Title            string   `json:"title"`
	} `json:"item"`
	Language []string `json:"language"`
	MimeType []string `json:"mime_type"`
	// Number         []string `json:"number"`
	NumberLccn     []string `json:"number_lccn"`
	OnlineFormat   []string `json:"online_format"`
	OriginalFormat []string `json:"original_format"`
	OtherTitle     []string `json:"other_title"`
	Partof         []string `json:"partof"`
	Resources      []struct {
		Files        int    `json:"files"`
		FulltextFile string `json:"fulltext_file"`
		Image        string `json:"image"`
		Pdf          string `json:"pdf"`
		Search       string `json:"search"`
		Segments     int    `json:"segments"`
		URL          string `json:"url"`
	} `json:"resources"`
	// Segments []struct {
	// 	Count int    `json:"count"`
	// 	Link  string `json:"link"`
	// 	URL   string `json:"url"`
	// } `json:"segments"`
	ShelfID string `json:"shelf_id"`
	// Site      []string  `json:"site"`
	Subject   []string  `json:"subject"`
	Timestamp time.Time `json:"timestamp"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
}

// Use the LCCN as a string representation of an item.
func (item ItemResult) String() string {
	return item.NumberLccn[0]
}

// Save serializes an item to the database.
func (item ItemResult) Save() error {
	query := `
	INSERT INTO items(id, lccn, url, date, subjects, title, api) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	// If we can convert the date to an integer do that, otherwise keep NULL
	var date sql.NullInt64
	parsed, err := strconv.Atoi(item.Date)
	if err == nil {
		date.Scan(parsed)
	}

	// Convert the rest of the data back to JSON to stuff into a DB column
	api, _ := json.Marshal(item)

	stmt, err := app.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(item.ID, item.NumberLccn[0], item.URL, date,
		item.Item.Subjects, item.Title, api)
	if err != nil {
		return err
	}

	return nil

}
