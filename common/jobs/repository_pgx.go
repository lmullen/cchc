package jobs

import (
	"context"

	"github.com/google/uuid"
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

// GetFullText finds a full text job by ID from the repository
func (r *Repo) GetFullText(ctx context.Context, id uuid.UUID) (*FullText, error) {
	query := `
	SELECT 
		id, item_id, has_ft_method, destination, started, finished
	FROM
		jobs.fulltext
	WHERE id = $1;
	`

	job := FullText{}

	err := r.db.QueryRow(ctx, query, id).Scan(&job.ID, &job.ItemID, &job.HasFTMethod,
		&job.Destination, &job.Started, &job.Finished)
	if err != nil {
		return nil, err
	}

	return &job, nil

}

// SaveFullText serializes a job to the database
func (r *Repo) SaveFullText(ctx context.Context, job *FullText) error {
	query := `
	INSERT INTO jobs.fulltext
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id) DO UPDATE
	SET
	item_id = $2,
	has_ft_method = $3,
	destination = $4,
	started = $5,
	finished = $6;
	`

	_, err := r.db.Exec(ctx, query,
		job.ID, job.ItemID, job.HasFTMethod, job.Destination,
		job.Started, job.Finished)
	if err != nil {
		return err
	}

	return nil

}
