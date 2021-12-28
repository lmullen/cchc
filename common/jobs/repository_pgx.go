package jobs

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
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

// CreateJobForUnqueued creates a job associated with a single item for a single
// destination. It should be safe for multiple workers to use.
//
// A potential problem is that occasionally a duplicate item key might be selected
// simultaneously by multiple workers. If that is the case, the database will still
// prevent the creation of a duplicate job for the same item and the same destination.
// This error could be safely ignored.
func (r *Repo) CreateJobForUnqueued(ctx context.Context, destination string) (*FullText, error) {

	query := `
	SELECT id
	FROM   items i
	WHERE  NOT EXISTS (
		SELECT item_id, destination
		FROM   jobs.fulltext
		WHERE  item_id=i.id AND destination = $1
	 )
	FOR NO KEY UPDATE OF i SKIP LOCKED
	LIMIT 1;
	`

	// query := `
	// SELECT i.id
	// FROM items i
	// LEFT JOIN jobs.fulltext j ON i.id = j.item_id
	// WHERE j.item_id IS NULL
	// FOR UPDATE OF i SKIP LOCKED
	// LIMIT 1;
	// `

	timeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// This operation has to happen in a transaction to be sure that the
	// `FOR UPDATE SKIP LOCKED` gets a unique item for each instance of this
	// function that might be running.
	tx, err := r.db.Begin(timeout)
	if err != nil {
		return nil, err
	}

	// Get a single item ID (guaranteeed to be unique, no matter how many instances
	// of this function are running) which we can use to make a job.
	var itemID string
	err = tx.QueryRow(timeout, query, destination).Scan(&itemID)
	if err != nil {
		tx.Rollback(timeout)
		if err == pgx.ErrNoRows {
			return nil, ErrAllQueued
		}
		return nil, err
	}

	// Now we have the item ID, so create a job for that item and that destination.
	job := NewFullText(itemID, destination)
	err = r.SaveFullText(timeout, job)
	if err != nil {
		tx.Rollback(timeout)
		return nil, err
	}

	// We saved the job so commit the transaction
	err = tx.Commit(timeout)
	if err != nil {
		// At this point, even if there is an error, we've saved the job, so return it
		return job, err
	}

	// Succesful so return
	return job, nil

}

// GetReadyJob gets a job for a particular destination that is available and marks it as started.
func (r *Repo) GetReadyJob(ctx context.Context, destination string) (*FullText, error) {
	timeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	getQuery := `
	SELECT id, item_id, destination, started, finished, status
	FROM jobs.fulltext
	WHERE status = 'ready' AND destination = $1
	FOR UPDATE SKIP LOCKED
	LIMIT 1;
	`

	// This operation has to happen in a transaction to be sure that the
	// `FOR UPDATE SKIP LOCKED` gets a unique item for each instance of this
	// function that might be running.
	tx, err := r.db.Begin(timeout)
	if err != nil {
		return nil, err
	}

	job := FullText{}

	err = tx.QueryRow(timeout, getQuery, destination).Scan(
		&job.ID,
		&job.ItemID,
		&job.Destination,
		&job.Started,
		&job.Finished,
		&job.Status,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNoJobs
		}
		return nil, err
	}

	updateQuery := `
	UPDATE jobs.fulltext 
	SET 
		status = 'running',
		started = $1
	WHERE id = $2;
	`

	now := time.Now()

	_, err = tx.Exec(ctx, updateQuery, now, job.ID)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	job.Started.Scan(now)
	job.Status = "running"

	err = tx.Commit(timeout)
	if err != nil {
		return &job, err
	}

	return &job, nil

}
