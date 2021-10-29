package results

import (
	"github.com/google/uuid"
)

// Quotation represents an instance of a biblical quotation in an item
type Quotation struct {
	JobID       uuid.UUID
	ItemID      string
	ReferenceID string
	VerseID     string
	Probability float64
}

// NewQuotation creates a new quotation object
func NewQuotation(JobID uuid.UUID, ItemID, ReferenceID, VerseID string, Probability float64) *Quotation {
	return &Quotation{
		JobID:       JobID,
		ItemID:      ItemID,
		ReferenceID: ReferenceID,
		VerseID:     VerseID,
		Probability: Probability,
	}

}
