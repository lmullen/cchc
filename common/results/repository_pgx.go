package results

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Repo is a data store using PostgreSQL with the pgx native interface.
type Repo struct {
	db *pgxpool.Pool
}

// NewRepo returns an item repo using PostgreSQL with the pgx native interface.
func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

// Save serializes a job to the database
func (r *Repo) Save(ctx context.Context, q *Quotation) error {
	query := `
	INSERT INTO results.biblical_quotations
	VALUES ($1, $2, $3, $4, $5);
	`

	_, err := r.db.Exec(ctx, query, q.JobID, q.ItemID, q.ReferenceID, q.VerseID, q.Probability)
	if err != nil {
		return err
	}

	return nil

}
