// Package results creates types for storing results from various programs in the database.
//
// It provides a Repository interface for generalized interactions storing and
// retrieving results from a data store, as well as a concrete type that implements
// that interface for a PostgreSQL database using the pgx package.
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
