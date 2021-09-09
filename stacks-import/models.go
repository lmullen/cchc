package main

import (
	"encoding/json"
	"time"
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
	PublicationDate string          `json:"publication_date"`
	Date            time.Time       `json:"date"`
	Year            int             `json:"year"`
	SubjectFull     []string        `json:"subject_full"`
	Subject         []string        `json:"subject"`
	Person          []string        `json:"person"`
	Place           json.RawMessage `json:"place"`
	Language        []string        `json:"language"`
	Form            json.RawMessage `json:"form"`
	File            json.RawMessage `json:"file"`
	Text            string          `json:"text_en"`
}
