package jobs

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// FulltextPredict is a record of how full text was passed on to the
// prediction model for each item.
type FulltextPredict struct {
	ID          uuid.UUID
	ItemID      string
	ResourceSeq sql.NullInt64
	FileSeq     sql.NullInt64
	FormatSeq   sql.NullInt64
	Level       sql.NullString
	Source      sql.NullString
	HasFTMethod bool
	Started     sql.NullTime
	Finished    sql.NullTime
}

// Repository is an interface describing a data store for jobs.
type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*FulltextPredict, error)
	Save(ctx context.Context, job *FulltextPredict) error
}
