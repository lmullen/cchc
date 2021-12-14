package jobs

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// Fulltext is a record of how jobs with the full text of items were passed on to
// some particular destination. It can handle recording full text jobs at either
// the item (i.e., resource) or file level, and it can handle multiple destinations.
type Fulltext struct {
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
	Queue       string
}

// Repository is an interface describing a data store for jobs.
type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*Fulltext, error)
	Save(ctx context.Context, job *Fulltext) error
}
