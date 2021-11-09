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

// Save serializes a job to the database
func (r *Repo) Save(ctx context.Context, job *Fulltext) error {
	query := `
	INSERT INTO jobs.fulltext_predict
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT (id) DO UPDATE
	SET 
	item_id = $2,
	resource_seq = $3,
	file_seq = $4,
	format_seq = $5,
	level = $6,
	source = $7,
	has_ft_method = $8,
	started = $9,
	finished = $10
	destination = $11
	;
	`

	_, err := r.db.Exec(ctx, query, job.ID, job.ItemID, job.ResourceSeq, job.FileSeq,
		job.FormatSeq, job.Level, job.Source, job.HasFTMethod, job.Started, job.Finished, job.Destination)
	if err != nil {
		return err
	}

	return nil

}

// Get finds a job by ID from the repository
func (r *Repo) Get(ctx context.Context, id uuid.UUID) (*Fulltext, error) {
	query := `
	SELECT 
		id, item_id, resource_seq, file_seq, format_seq, level, source, has_ft_method, started, finished, destination
	FROM
		jobs.fulltext_predict
	WHERE id = $1;
	`

	job := Fulltext{}

	err := r.db.QueryRow(ctx, query, id).Scan(&job.ID, &job.ItemID, &job.ResourceSeq,
		&job.FileSeq, &job.FormatSeq, &job.Level, &job.Source, &job.HasFTMethod,
		&job.Started, &job.Finished, &job.Destination)
	if err != nil {
		return nil, err
	}

	return &job, nil

}
