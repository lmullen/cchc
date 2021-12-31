package results

import (
	"context"

	"github.com/google/uuid"
)

// Repository is an interface describing a data store
type Repository interface {
	SaveQuotation(ctx context.Context, q *Quotation) error
	SaveLanguages(ctx context.Context, jobID uuid.UUID, itemID string, languages map[string]int) error
}
