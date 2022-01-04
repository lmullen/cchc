package main

import (
	"context"

	"github.com/lmullen/cchc/common/jobs"

	log "github.com/sirupsen/logrus"
)

func setJobs(jobs []*jobs.FullText, succeed bool) {
	// Set the job status
	for _, job := range jobs {
		if succeed {
			job.Finish()
		} else {
			job.Fail()
		}
	}

	// Save the jobs to the database
	for _, job := range jobs {
		err := app.JobsRepo.SaveFullText(context.TODO(), job)
		if err != nil {
			log.WithError(err).WithField("job", job).Error("Error saving job status")
		}
	}

}
