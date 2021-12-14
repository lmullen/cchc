package main

import (
	"context"
	"fmt"

	"github.com/lmullen/cchc/common/jobs"
	"github.com/lmullen/cchc/common/messages"
)

// EnqueueItemFulltext determines whether an item has full text available,
// creates jobs to keep track of the pieces of text, and sends the text to the
// appropriate message queue.
//
// The messages repository controls which queue the job is sent to.
func EnqueueItemFulltext(ctx context.Context, itemID string, mr messages.Repository) error {
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
				job := &jobs.Fulltext{}
				job.Create(item.ID, mr.QueueName(), true)
				fulltext := job.PlainTextFullText(file)
				fulltext = app.stripXML.Sanitize(fulltext)
				msg := messages.NewFullTextMsg(job.ID, job.ItemID, fulltext)

				err := mr.Send(ctx, msg)
				if err != nil {
					return fmt.Errorf("Error sending job to message queue: %w", err)
				}
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
				job := &jobs.Fulltext{}
				job.Create(item.ID, mr.QueueName(), true)
				fulltext := job.XMLFullText(file)
				fulltext = app.stripXML.Sanitize(fulltext)
				msg := messages.NewFullTextMsg(job.ID, job.ItemID, fulltext)

				err := mr.Send(ctx, msg)
				if err != nil {
					return fmt.Errorf("Error sending job to message queue: %w", err)
				}
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
	// if !hasFulltext {
	// job := &jobs.FulltextPredict{}
	// job.Create(item.ID, false)
	// 	err = app.JR.Save(ctx, job)
	// 	if err != nil {
	// 		return fmt.Errorf("Error saving job: %w", err
	// 	}
	// }

	return nil

}
