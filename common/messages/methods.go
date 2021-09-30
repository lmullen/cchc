package messages

import "github.com/google/uuid"

// NewFullTextMsg creates a pointer to a new FullTextPredict message
func NewFullTextMsg(JobID uuid.UUID, text string) *FullTextPredict {
	return &FullTextPredict{
		JobID:    JobID,
		FullText: text,
	}
}
