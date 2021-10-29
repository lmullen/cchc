package results

import "context"

// Repository is an interface describing a data store
type Repository interface {
	Save(ctx context.Context, q *Quotation) error
}
