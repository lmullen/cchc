package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Item is a representation of an item in the LOC collection returned from the API.
type Item struct {
	ID              string
	URL             string
	Title           string
	Date            string
	Subjects        []string
	Fulltext        string
	FulltextService string
	FulltextFile    string
	Timestamp       int64
	API             []byte // The entire API response stored as JSONB
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
		ID       string   `json:"id"`
		URL      string   `json:"url"`
		Date     string   `json:"date"`
		Subjects []string `json:"subject_headings"`
		Title    string   `json:"title"`
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

// Exists checks whether an item has been saved to the database or not.
func (i Item) Exists() (bool, error) {
	var exists bool
	err := app.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM items WHERE id=$1)", i.ID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, fmt.Errorf("Error checking if item %s exists: %w", i.ID, err)
	}
	return exists, nil
}

// Save serializes an item to the database
func (i Item) Save() error {
	query := `
	INSERT INTO items (id, url, title, date, subjects, 
		                fulltext, fulltext_service, fulltext_file,
										timestamp, api)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	ON CONFLICT (id) DO UPDATE
	SET
	  url              = $2,
		title            = $3,
		date             = $4,
		subjects         = $5,
		fulltext         = $6,
		fulltext_service = $7,
		fulltext_file    = $8,
		timestamp        = $9,
		api              = $10;
	`

	_, err := app.DB.Exec(query, i.ID, i.URL, i.Title, i.Date, i.Subjects,
		i.Fulltext, i.FulltextService, i.FulltextFile, i.Timestamp, i.API)
	fmt.Println(err)
	if err != nil {
		return fmt.Errorf("Error saving item %s to database: %w", i, err)
	}

	return nil
}

// String returns a string representation of an item.
func (i Item) String() string {
	return i.ID
}

// Fetch gets an item's metadata from the LOC.gov API.
func (i *Item) Fetch() error {

	app.Limiters.Items.Take()

	u, _ := url.Parse(i.URL)
	remove := []string{"more_like_this", "related_items", "cite_this"}
	options := url.Values{
		"at!": []string{strings.Join(remove, ",")},
		"fo":  []string{"json"},
	}
	u.RawQuery = options.Encode()
	url := u.String()

	log.WithField("url", url).Info("Fetching item metadata")

	response, err := app.Client.Get(url)
	if err != nil {
		return fmt.Errorf("Error fetching item %s: %w", url, err)
	}

	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"http_error": response.Status,
			"http_code":  response.StatusCode,
			"url":        url,
		}).Warn("HTTP error when fetching from API")
		quitIfBlocked(response.StatusCode)
		return fmt.Errorf("HTTP error: %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error reading HTTP response body: %w", err)
	}

	var result ItemResponse

	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("Error unmarshalling item metadata: %w", err)
	}

	i.ID = result.ItemDetails.ID
	i.URL = result.ItemDetails.URL
	i.Title = result.ItemDetails.Title
	i.Date = result.ItemDetails.Date
	i.Subjects = result.ItemDetails.Subjects
	i.Timestamp = result.Timestamp
	i.API = data

	// TODO Getting the full text fields here is janky. Not sure how consistent
	// the API is.
	for _, v := range result.Resources[0].Files[0] {
		if v.Fulltext != "" {
			i.Fulltext = v.Fulltext
		}
		if v.FulltextService != "" {
			i.FulltextService = v.FulltextService
		}
	}
	i.FulltextFile = result.Resources[0].FulltextFile

	return nil

}
