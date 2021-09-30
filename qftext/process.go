package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/lmullen/cchc/common/jobs"
	"github.com/lmullen/cchc/common/messages"
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

	// FULL TEXT CHECK 1: Has text/plain mimetype with fulltext field
	if !hasFulltext {
		for _, file := range item.Files {
			if file.Mimetype.Valid && file.Mimetype.String == "text/plain" && file.FullText.Valid {
				job := &jobs.FulltextPredict{}
				job.Create(item.ID, true)
				fulltext := job.PlainTextFullText(file)
				fulltext = app.stripXML.Sanitize(fulltext)
				msg := messages.NewFullTextMsg(job.ID, fulltext)

				// SEND MESSAGE HERE
				log.Println("Plain text", msg.JobID)
				err = app.JR.Save(ctx, job)
				if err != nil {
					return fmt.Errorf("Error saving job: %w", err)
				}
				hasFulltext = true // Keep track that we found full text
			}
		}
	}

	// FULL TEXT CHECK 1: Has text/xml mimetype with fulltext field
	if !hasFulltext {
		for _, file := range item.Files {
			if file.Mimetype.Valid && file.Mimetype.String == "text/xml" && file.FullText.Valid {
				job := &jobs.FulltextPredict{}
				job.Create(item.ID, true)
				fulltext := job.XMLFullText(file)
				fulltext = app.stripXML.Sanitize(fulltext)
				msg := messages.NewFullTextMsg(job.ID, fulltext)

				// SEND MESSAGE HERE
				log.Println("XML", msg.JobID)
				err = app.JR.Save(ctx, job)
				if err != nil {
					return fmt.Errorf("Error saving job: %w", err)
				}
				hasFulltext = true // Keep track that we found full text
			}
		}
	}

	// These jobs have no method for finding the full text, so create a job to
	// skip them
	if !hasFulltext {
		job := &jobs.FulltextPredict{}
		job.Create(item.ID, false)
		err = app.JR.Save(ctx, job)
		if err != nil {
			return fmt.Errorf("Error saving job: %w", err)
		}
	}

	return nil

}
