package main

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/lmullen/cchc/common/jobs"
	log "github.com/sirupsen/logrus"
)

// ProcessUnqueued iterates through the items which have full text and which do
// not have a job in the table. It then creates a new job for each of those
// items.
func ProcessUnqueued(ctx context.Context) error {

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
			return err
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
// job to keep track of whether
func ProcessItem(ctx context.Context, itemID string) error {
	// log.WithField("item-id", itemID).Debug("Checking item's full text for processing")

	item, err := app.IR.Get(ctx, itemID)
	if err != nil {
		return err
	}

	// FULL TEXT CHECK 1: Has plain text files with full text service
	job := &jobs.JobFulltextPredict{}
	err = job.Create(item.ID)
	if err != nil {
		return err
	}

	for _, file := range item.Files {
		if file.Mimetype.Valid && file.Mimetype.String == "text/plain" && file.FullTextService.Valid {
			job.HasFTMethod = true
			job.Level.Scan("file")
			job.Source.Scan("plain text full service")
			SendFullText(job.ID.String(), file.FullTextService.String)
		}
	}

	err = app.JR.Save(ctx, job)
	if err != nil {
		return err
	}
	log.Info(job)

	return nil

}
