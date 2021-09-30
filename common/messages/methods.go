package messages

import "github.com/google/uuid"

// NewFullTextMsg creates a pointer to a new FullTextPredict message
func NewFullTextMsg(job uuid.UUID, item string, text string) *FullTextPredict {
	return &FullTextPredict{
		JobID:    job,
		ItemID:   item,
		FullText: text,
	}
}
