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
		id, item_id, destination, started, finished, status
	FROM
		jobs.fulltext
	WHERE id = $1;
	`

	job := FullText{}

	err := r.db.QueryRow(ctx, query, id).Scan(&job.ID, &job.ItemID,
		&job.Destination, &job.Started, &job.Finished, &job.Status)
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
	destination = $3,
	started = $4,
	finished = $5,
	status = $6;
	`

	_, err := r.db.Exec(ctx, query,
		job.ID, job.ItemID, job.Destination, job.Started, job.Finished, job.Status)
	if err != nil {
		return err
	}

	return nil

}

// func (r *Repo) CreateJobForUnqueued(ctx context.Context, destination string) (*FullText, error) {

// 	query := `
// 	SELECT id
// 	FROM   items i
// 	WHERE  NOT EXISTS (
// 		SELECT item_id, destination
// 		FROM   jobs.fulltext
// 		WHERE  item_id=i.id AND destination = $1
//    )
// 	FOR UPDATE SKIP LOCKED
// 	LIMIT 1;
// 	`

// 	timeout, cancel := context.WithTimeout(ctx, 60*time.Minute)
// 	defer cancel()

// 	// This operation has to happene in a transaction to be sure that the
// 	// `FOR UPDATE SKIP LOCKED` gets a unique item for each instance of this
// 	// function that might be running.
// 	tx, err := r.db.Begin(timeout)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get a single item ID (guaranteeed to be unique, no matter how many instances
// 	// of this function are running) which we can use to make a job.
// 	var itemID string
// 	err = tx.QueryRow(timeout, query, destination).Scan(&itemID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Now we have the item ID, so create a job for that item and that destination.

// }
