package main

import (
	"os"
	"os/exec"
	"time"

	"github.com/lmullen/cchc/common/jobs"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func processBatchOfDocs(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		timeout, cancel := context.WithTimeout(ctx, jobtimeout)
		defer cancel()

		// Keep track of the jobs in this batch
		jobsInBatch := make([]*jobs.FullText, 0, itemsPerBatch)
		// Keep track of the specific pages in this batch
		docsInBatch := make([]*Doc, 0, itemsPerBatch)

		log.Debugf("Collecting a batch of up to %v items or %v pages", itemsPerBatch, pagesPerBatch)

		// We are going to get as many items as we want to form a batch.
		// Keep in mind that each item may contain many subdocuments.
		// We will skip jobs if they don't have full text. And if we run to the end
		// the available batches, we will start early with what we have.
		for len(jobsInBatch) < itemsPerBatch && len(docsInBatch) < pagesPerBatch {
			job, err := app.JobsRepo.GetReadyJob(timeout, queue)
			if err != nil {
				if err == jobs.ErrNoJobs && len(jobsInBatch) > 0 {
					// There are no more jobs, but we already have at least one job,
					// so break out of this loop and process what we've got.
					break
				}
				if err == jobs.ErrNoJobs && len(jobsInBatch) == 0 {
					// There are no more jobs and we have no jobs at all. So wait to check
					// if there are more jobs to process.
					log.Infof("No ready jobs, so waiting %s to check again", waittime)
					select {
					case <-ctx.Done():
						return
					case <-time.After(waittime):
						log.Info("Checking again whether there are jobs to be processed")
						return
					}
				}
				// We have an error but it is unanticipated so log it and quit this function
				log.WithError(err).Error("Error getting a job that is ready")
				return
			}

			// Start this job so nobody else gets it
			job.Start()
			err = app.JobsRepo.SaveFullText(timeout, job)
			if err != nil {
				log.WithError(err).WithField("job", job).Error("Error saving job status")
				continue
			}

			// Get the item
			item, err := app.ItemsRepo.Get(timeout, job.ItemID)
			if err != nil {
				log.WithError(err).WithField("job", job).Error("Error getting item for job")
				job.Fail()
				err = app.JobsRepo.SaveFullText(timeout, job)
				if err != nil {
					log.WithError(err).WithField("job", job).Error("Error saving job status")
				}
				continue
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
				continue // Keep going collecting items
			}

			// Keep track of the fact that we are running this job
			jobsInBatch = append(jobsInBatch, job)

			// Add the pages to the batch
			for _, page := range pages {
				doc := NewDoc(job, item, page)
				docsInBatch = append(docsInBatch, doc)
			}
		}

		// At this point we have a complement of jobs, so start the process of running them
		log.Debugf("Running quotation finder on a batch of %v items and %v pages", len(jobsInBatch), len(docsInBatch))

		// Write the full text to a temporary CSV
		docsFile, err := writeDocsCSV(docsInBatch)
		if err != nil {
			log.WithError(err).Error("Error writing CSV to send to prediction model")
		}

		// Create a temp file for output.
		predictionsFile, err := os.CreateTemp("", "prediction-*.csv")
		if err != nil {
			log.WithError(err).Error("Error creating temporary file for predictions")
		}
		predictionsFile.Close()

		cmd := exec.CommandContext(timeout,
			"Rscript", "/predictor/id-quotations.R",
			"--bible", "bible-payload.rda",
			"--model", "prediction-payload.rda",
			"--verbose", "2",
			"--out", predictionsFile.Name(),
			// "--potential",
			docsFile,
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.WithError(err).WithField("R-output", string(output)).Error("Problem running prediction model in R")
			setJobs(jobsInBatch, false)
			return
		}

		// Get the predictions back from a temporary file and write them to the database
		err = processPredictionsCSV(ctx, predictionsFile.Name())
		if err != nil {
			log.WithError(err).Error("Error getting results from prediction model")
			setJobs(jobsInBatch, false)
			return
		}

		setJobs(jobsInBatch, true)

		// Clean up the temporary files
		err = os.Remove(predictionsFile.Name())
		if err != nil {
			log.WithError(err).Warn("Problem removing the temporary files")
		}
		err = os.Remove(docsFile)
		if err != nil {
			log.WithError(err).Warn("Problem removing the temporary files")
		}

		log.Debugf("Finished processing a batch of %v items and %v pages", len(jobsInBatch), len(docsInBatch))

	}
}
