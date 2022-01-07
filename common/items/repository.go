// Package items allows for interactions with items in the Library of Congress digital collections.
//
// It provides a Repository interface for generalized interactions storing and
// retrieving items from a data store, as well as a concrete type that implements
// that interface for a PostgreSQL database using the pgx package.
package items

import "context"

// Repository is an interface describing a data store for items.
type Repository interface {
	Get(ctx context.Context, ID string) (*Item, error)
	GetAllUnfetched(ctx context.Context) ([]string, error)
	Save(ctx context.Context, item *Item) error
}
