package messages

import (
	"github.com/google/uuid"
)

// FullTextPredict represents a job for a file or resource passed as a message
// to the predictor.
type FullTextPredict struct {
	JobID    uuid.UUID `json:"job_id"`
	ItemID   string    `json:"item_id"`
	FullText string    `json:"full_text"`
}
