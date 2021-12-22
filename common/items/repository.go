package items

import "context"

// Repository is an interface describing a data store for items.
type Repository interface {
	Get(ctx context.Context, ID string) (*Item, error)
	GetAllUnfetched(ctx context.Context) ([]string, error)
	Save(ctx context.Context, item *Item) error
}
