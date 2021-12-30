package main

import (
	"context"

	"github.com/lmullen/cchc/common/jobs"
	log "github.com/sirupsen/logrus"
)

func processDocument(job *jobs.FullText) error {

	// We've received a job. We need to calculate the language stats if possible,
	// save the results to the database, and if not, then skip the job.

	// This is a job we are processing so fail if it takes too long.
	ctx, cancel := context.WithTimeout(context.Background(), jobtimeout)
	defer cancel()

	// The job has the item we need, so get the item from the database.
	item, err := app.ItemsRepo.Get(ctx, job.ItemID)
	if err != nil {
		job.Fail()
		errSave := app.JobsRepo.SaveFullText(ctx, job)
		if errSave != nil {
			log.WithError(err).WithField("job", job).Error("Error saving failed job status")
		}
		return err
	}

	// Get the full text for the item
	pages, has := item.FullText()

	// It's possible we don't have full text. If so, skip the job.
	if !has {
		job.Skip()
		errSave := app.JobsRepo.SaveFullText(ctx, job)
		if errSave != nil {
			log.WithError(err).WithField("job", job).Error("Error saving job status")
		}
		return nil
	}

	// Make a results map that will be shared for all pages in the item.
	results := make(LanguageStats)

	for _, p := range pages {
		err = CalculateLanguages(p.Text, results)
		if err != nil {
			job.Fail()
			errSave := app.JobsRepo.SaveFullText(ctx, job)
			if errSave != nil {
				log.WithError(err).WithField("job", job).Error("Error saving job status")
			}
			return err
		}
	}

	// This is where we should save the results to the database. But for now, just print them.
	log.WithField("job", job).WithField("results", results).Debug("Successfully processed job")

	// The job was successful
	job.Finish()
	errSave := app.JobsRepo.SaveFullText(ctx, job)
	if errSave != nil {
		log.WithError(err).WithField("job", job).Error("Error saving job status")
	}

	return nil

}
