package messages

import (
	"context"
	"database/sql"
)

// FullTextPredict represents a job for a file or resource passed as a message
// to the predictor.
type FullTextPredict struct {
	JobID       string        `json:"job_id"`
	ItemID      string        `json:"item_id"`
	ResourceSeq sql.NullInt64 `json:"resource_seq"`
	FileSeq     sql.NullInt64 `json:"file_seq"`
	FormatSeq   sql.NullInt64 `json:"format_seq"`
	FullText    string        `json:"full_text"`
}

// Repository is an interface describing a datastore for messages
type Repository interface {
	Send(ctx context.Context, text *FullTextPredict) error
	// Receive(ctx context.Context, text *FullTextPredict) error
}
