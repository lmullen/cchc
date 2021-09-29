package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/lmullen/cchc/common/jobs"
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

// ProcessItem determines whether an item has full text available, and creates a
// job to keep track of whether its text has been sent to the prediction model.
func ProcessItem(ctx context.Context, itemID string) error {
	item, err := app.IR.Get(ctx, itemID)
	if err != nil {
		return err
	}

	// Keep track of whether we have found full text yet
	hasFulltext := false

	// FULL TEXT CHECK 1: Has plain text files with full text service
	if !hasFulltext {
		for _, file := range item.Files {
			if file.Mimetype.Valid && file.Mimetype.String == "text/plain" && file.FullTextService.Valid {
				job := &jobs.FulltextPredict{}
				job.Create(item.ID)
				job.FulltextFromFile(file, "plain text/full text service")
				job.Start()
				SendFullText(job.ID.String(), file.FullTextService.String)
				hasFulltext = true // Keep track that we found full text
				err = app.JR.Save(ctx, job)
				if err != nil {
					return fmt.Errorf("Error saving job: %w", err)
				}
			}
		}

	}

	if !hasFulltext {
		job := &jobs.FulltextPredict{}
		job.Create(item.ID)
		if err != nil {
			return fmt.Errorf("Error creating item: %w", err)
		}
		err = app.JR.Save(ctx, job)
		if err != nil {
			log.Println(job)
			return fmt.Errorf("Error saving job: %w", err)
		}
		log.WithField("job", job).Info("Skipped because no full text")
	}

	if err != nil {
		return err
	}

	return nil

}
