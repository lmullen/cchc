package main

import (
	"github.com/google/uuid"
	"github.com/lmullen/cchc/common/items"
	"github.com/lmullen/cchc/common/jobs"
)

// NewDoc creates a document from a job, item, and page of an item
func NewDoc(job *jobs.FullText, item *items.Item, text items.PlainText) *Doc {
	return &Doc{
		JobID:  job.ID,
		ItemID: item.ID,
		Text:   text.Text,
	}
}

// Doc is used by the CSV writer and quotation finder.
type Doc struct {
	JobID  uuid.UUID
	ItemID string
	Text   string
}

// CSVRow converts a Doc into a format for writing to a CSV.
func (doc *Doc) CSVRow() []string {
	out := make([]string, 3)
	out[0] = doc.JobID.String()
	out[1] = doc.ItemID
	out[2] = doc.Text
	return out
}
