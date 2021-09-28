package jobs

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// JobFulltextPredict is a record of how full text was passed on to the
// prediction model for each item.
type JobFulltextPredict struct {
	ID          uuid.UUID
	ItemID      string
	Level       sql.NullString
	Source      sql.NullString
	HasFTMethod bool
	Started     sql.NullTime
	Finished    sql.NullTime
}

// Repository is an interface describing a data store for jobs.
type Repository interface {
	Save(ctx context.Context, job *JobFulltextPredict) error
}
