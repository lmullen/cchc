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

// CSVRow converts a FullTextPredict message into a format for writing to a CSV.
func (f *FullTextPredict) CSVRow() []string {
	out := make([]string, 3)
	out[0] = f.JobID.String()
	out[1] = f.ItemID
	out[2] = f.FullText
	return out
}
