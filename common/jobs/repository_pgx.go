package jobs

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Repo is a data store using PostgreSQL with the pgx native interface.
type Repo struct {
	db *pgxpool.Pool
}

// NewJobsRepo returns an item repo using PostgreSQL with the pgx native interface.
func NewJobsRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

// Save serializes a job to the database
func (r *Repo) Save(ctx context.Context, job *FulltextPredict) error {
	query := `
	INSERT INTO jobs.fulltext_predict
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	_, err := r.db.Exec(ctx, query, job.ID, job.ItemID, job.ResourceSeq, job.FileSeq,
		job.FormatSeq, job.Level, job.Source, job.HasFTMethod, job.Started, job.Finished)
	if err != nil {
		return err
	}

	return nil

}
