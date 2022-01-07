// Package jobs keeps track of work to be done against the database, including queuing
//
// It provides a Repository interface for generalized interactions storing and
// retrieving jobs from a data store, as well as a concrete type that implements
// that interface for a PostgreSQL database using the pgx package.
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
	GetReadyJob(ctx context.Context, destination string) (*FullText, error)
}
