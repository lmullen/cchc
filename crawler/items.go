package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Item is a representation of an item in the LOC collection returned from the API.
type Item struct {
	ID        string
	URL       string
	Title     string
	Year      sql.NullInt32
	Date      string
	Subjects  []string
	Timestamp int64
	Resources []ItemResource
	Files     []ItemFile
	API       []byte // The entire API response stored as JSONB
}

// ItemResource contains a description of some kind of source.
type ItemResource struct {
	ItemID       string
	ResourceSeq  int
	FullTextFile sql.NullString
	DJVUTextFile sql.NullString
	Image        sql.NullString
	PDF          sql.NullString
	URL          sql.NullString
	Caption      sql.NullString
}

// ItemFile contains a file pointing to some kind of resource. Unlike the
// LOC.gov API, it does not make a firm distinction between a file and format.
type ItemFile struct {
	ItemID          string
	ResourceSeq     int
	FileSeq         int
	FormatSeq       int
	Mimetype        sql.NullString
	FullText        sql.NullString
	FullTextService sql.NullString
	WordCoordinates sql.NullString
	URL             sql.NullString
	Info            sql.NullString
	Use             sql.NullString
}

// ItemResponse represents an item-level object returned from the API. Many more fields
// are returned and will be stored in the database as a JSONB field, but these
// are the ones that will be serialized to regular database fields.
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
		FulltextFile string `json:"fulltext_file,omitempty"`
		DJVUTextFile string `json:"djvu_text_file,omitempty"`
		Image        string `json:"image,omitempty"`
		PDF          string `json:"pdf,omitempty"`
		URL          string `json:"url,omitempty"`
		Caption      string `json:"caption,omitempty"`
		Files        [][]struct {
			Mimetype        string `json:"mimetype,omitempty"`
			Fulltext        string `json:"fulltext,omitempty"`
			FulltextService string `json:"fulltext_service,omitempty"`
			WordCoordinates string `json:"word_coordinates,omitempty"`
			URL             string `json:"url,omitempty"`
			Info            string `json:"info,omitempty"`
			Use             string `json:"use,omitempty"`
		} `json:"files"`
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

// Fetched checks whether an item's metadata has already been fetched
// from the API. This will also return false if the item has not been saved at
// all.
func (i Item) Fetched() (bool, error) {
	var fetched bool
	query := `SELECT EXISTS (SELECT 1 FROM items WHERE id=$1 AND api IS NOT NULL)`
	err := app.DB.QueryRow(query, i.ID).Scan(&fetched)
	if err != nil && err != sql.ErrNoRows {
		return fetched, fmt.Errorf("Error checking if item %s has been fetched: %w", i.ID, err)
	}
	return fetched, nil
}

// Save serializes an item to the database
func (i Item) Save() error {
	itemQuery := `
	INSERT INTO items (id, url, title, year, date, subjects, timestamp, api)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (id) DO UPDATE
	SET
	  url              = $2,
		title            = $3,
		year             = $4,
		date             = $5,
		subjects         = $6,
		timestamp        = $7,
		api              = $8;
	`

	resourceQuery := `
	INSERT INTO resources (item_id, resource_seq, fulltext_file, djvu_text_file,
		image, pdf, url, caption)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	fileQuery := `
	INSERT INTO files (item_id, resource_seq, file_seq, format_seq,
	                   mimetype, fulltext, fulltext_service, word_coordinates,
										 url, info, use)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	// Use a transaction since we are writing to three tables
	tx, err := app.DB.Begin()
	if err != nil {
		return fmt.Errorf("Error creating transaction in database: %w", err)
	}

	_, err = tx.Exec(itemQuery, i.ID, i.URL, i.Title, i.Year, i.Date, i.Subjects,
		i.Timestamp, i.API)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error saving item %s to database: %w", i, err)
	}

	for _, r := range i.Resources {
		_, err = tx.Exec(resourceQuery, r.ItemID, r.ResourceSeq, r.FullTextFile,
			r.DJVUTextFile, r.Image, r.PDF, r.URL, r.Caption)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error saving item %s to database: %w", i, err)
		}
	}

	for _, f := range i.Files {
		_, err = tx.Exec(fileQuery, f.ItemID, f.ResourceSeq, f.FileSeq, f.FormatSeq,
			f.Mimetype, f.FullText, f.FullTextService, f.WordCoordinates, f.URL,
			f.Info, f.Use)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Error saving item %s to database: %w", i, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Error committing item to database: %w", err)
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

	log.WithField("url", url).Debug("Fetching item metadata")

	response, err := app.Client.Get(url)
	if err != nil {
		return fmt.Errorf("Error getting item over HTTP: %w", err)
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
		return fmt.Errorf("Error reading HTTP response body while fetching item: %w", err)
	}

	var result ItemResponse

	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("Error unmarshalling item metadata: %w", err)
	}

	i.ID = result.ItemDetails.ID
	i.URL = result.ItemDetails.URL
	i.Title = result.ItemDetails.Title
	i.Year = year(result.ItemDetails.Date)
	i.Date = result.ItemDetails.Date
	i.Subjects = result.ItemDetails.Subjects
	i.Timestamp = result.Timestamp
	i.API = data

	var temp [][]string

	for i, v := range temp {
		for j, w := range v {
			fmt.Println(i, v, j, w)

		}

	}

	// Iterate through all the files and formats to get the full text representations
	for resourceSeq, resource := range result.Resources {
		var r ItemResource
		r.ItemID = i.ID
		r.ResourceSeq = resourceSeq
		if resource.FulltextFile != "" {
			r.FullTextFile.Scan(resource.FulltextFile)
		}
		if resource.DJVUTextFile != "" {
			r.DJVUTextFile.Scan(resource.DJVUTextFile)
		}
		if resource.Image != "" {
			r.Image.Scan(resource.Image)
		}
		if resource.PDF != "" {
			r.PDF.Scan(resource.PDF)
		}
		if resource.URL != "" {
			r.URL.Scan(resource.URL)
		}
		if resource.Caption != "" {
			r.Caption.Scan(resource.Caption)
		}
		i.Resources = append(i.Resources, r)
		for fileSeq, file := range resource.Files {
			for formatSeq, format := range file {
				var f ItemFile
				f.ItemID = i.ID
				f.ResourceSeq = resourceSeq
				f.FileSeq = fileSeq
				f.FormatSeq = formatSeq
				if format.Mimetype != "" {
					f.Mimetype.Scan(format.Mimetype)
				}
				if format.Fulltext != "" {
					f.FullText.Scan(format.Fulltext)
				}
				if format.FulltextService != "" {
					f.FullTextService.Scan(format.FulltextService)
				}
				if format.WordCoordinates != "" {
					f.WordCoordinates.Scan(format.WordCoordinates)
				}
				if format.URL != "" {
					f.URL.Scan(format.URL)
				}
				if format.Info != "" {
					f.Info.Scan(format.Info)
				}
				if format.Use != "" {
					f.Use.Scan(format.Use)
				}
				i.Files = append(i.Files, f)
			}
		}
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
func (i Item) Msg() ItemMetadataMsg {
	return ItemMetadataMsg{
		ID:  i.ID,
		URL: i.URL,
	}
}

// EnqueueMetadata sends a message to the message queue to so that the item's
// metadata will put into a queue to be fetched.
func (i Item) EnqueueMetadata() error {

	json, err := json.Marshal(i.Msg())
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

// ProcessItemMetadata reads an item from the queue, fetches its metadata, and
// saves it to the database.
func ProcessItemMetadata(msg amqp.Delivery) {
	var item Item
	err := json.Unmarshal(msg.Body, &item)
	if err != nil {
		msg.Reject(false)
		log.WithError(err).WithField("msg", msg).Error("Failed to read body of message from queue")
		return
	}
	if fetched, err := item.Fetched(); fetched && err == nil {
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
