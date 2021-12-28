package jobs

import (
	"context"

	"github.com/google/uuid"
)

// Repository is an interface describing a data store for jobs.
type Repository interface {
	GetFullText(ctx context.Context, id uuid.UUID) (*FullText, error)
	SaveFullText(ctx context.Context, job *FullText) error
	CreateJobForUnqueued(ctx context.Context, destination string) (*FullText, error)
}
