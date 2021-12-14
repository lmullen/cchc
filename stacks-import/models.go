package main

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"
)

// Book represents a JSON object in the Stacks export, with metadata and full text
type Book struct {
	ItemType        json.RawMessage `json:"item_type"`
	LCCN            string          `json:"lccn"`
	ISBN            []string        `json:"isbn"`
	Title           string          `json:"title"`
	SortTitle       json.RawMessage `json:"sort_title"`
	Publisher       string          `json:"publisher"`
	Published       json.RawMessage `json:"published"`
	PublicationDate json.RawMessage `json:"publication_date"`
	Year            int             `json:"-"`
	SubjectFull     []string        `json:"subject_full"`
	Subject         []string        `json:"subject"`
	Person          []string        `json:"person"`
	Place           json.RawMessage `json:"place"`
	Language        []string        `json:"language"`
	Form            json.RawMessage `json:"form"`
	File            json.RawMessage `json:"file"`
	Text            string          `json:"text_en"`
}

func (b Book) String() string {
	return b.LCCN
}

// Exists checks whether a book as already been serialized to the database
func (b Book) Exists() (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM stacks_books WHERE lccn=$1);`
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var exists bool
	err := db.QueryRow(ctx, query, b.LCCN).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

// Save serializes a book to the database
func (b Book) Save() error {
	query := `
	INSERT INTO stacks_books(
		lccn, 
		isbn,
		title,
		publisher,
		year,
		subject_full,
		subject,
		person, 
		lang,
		text) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	ON CONFLICT DO NOTHING;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Exec(ctx, query, b.LCCN, b.ISBN, b.Title, b.Publisher,
		b.Year, b.SubjectFull, b.Subject, b.Person, b.Language, b.Text)

	return err

}

// SaveJSON adds the original metadata as JSON to the database
func (b Book) SaveJSON() error {
	query := `UPDATE stacks_books SET original_metadata=$1 WHERE lccn=$2`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	json, err := json.Marshal(b)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query, json, b.LCCN)

	return err
}
