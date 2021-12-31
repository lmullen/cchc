package results

import (
	"context"

	"github.com/google/uuid"
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

// SaveQuotation serializes a job to the database
func (r *Repo) SaveQuotation(ctx context.Context, q *Quotation) error {
	query := `
	INSERT INTO results.biblical_quotations (job_id, item_id, reference_id, verse_id, probability)
	VALUES ($1, $2, $3, $4, $5);
	`

	_, err := r.db.Exec(ctx, query, q.JobID, q.ItemID, q.ReferenceID, q.VerseID, q.Probability)
	if err != nil {
		return err
	}

	return nil

}

// SaveLanguages serializes the results of calculating languages to the database.
func (r *Repo) SaveLanguages(ctx context.Context, jobID uuid.UUID, itemID string, languages map[string]int) error {

	insert := `
			INSERT INTO results.languages (job_id, item_id, lang, sentences)
			VALUES ($1, $2, $3, $4);
		`

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Roll back the transaction if something goes wrong

	for lang, sent := range languages {
		_, err := tx.Exec(ctx, insert, jobID, itemID, lang, sent)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil

}
