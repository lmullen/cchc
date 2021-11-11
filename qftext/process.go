package main

import (
	"context"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

// FindUnprocessedItems iterates through the items which have full text and which do
// not have a job in the table. It then creates a new job for each of those
// items.
func FindUnprocessedItems(ctx context.Context) error {

	tx, err := app.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DECLARE unqcurs NO SCROLL CURSOR FOR SELECT id FROM jobs.fulltext_unqueued;`)
	if err != nil {
		return err
	}

	var item string
	for {
		err = tx.QueryRow(ctx, `FETCH NEXT FROM unqcurs;`).Scan(&item)
		if err != nil {
			if err == pgx.ErrNoRows {
				break
			}
			return err
		}
		err = ProcessItem(ctx, item)
		if err != nil {
			log.WithField("item", item).WithError(err).Error("Error processing item")
		}
	}

	_, err = tx.Exec(ctx, "CLOSE unqcurs;")
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil

}
