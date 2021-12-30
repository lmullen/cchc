package main

import (
	"context"
	"sync"
	"time"

	"github.com/lmullen/cchc/common/jobs"
	log "github.com/sirupsen/logrus"
)

func processJobs(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info("Checking whether there are jobs to be processed")
	for {
		select {
		case <-ctx.Done():
			log.Info("Stopped processing jobs")
			return
		default:
			timeoutGet, cancelGet := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancelGet()
			job, err := app.JobsRepo.GetReadyJob(timeoutGet, queue)
			if err != nil {
				if err == jobs.ErrNoJobs {
					log.Infof("No ready jobs, so waiting %s to check again", waittime)
					select {
					case <-ctx.Done():
						return
					case <-time.After(waittime):
						log.Info("Checking again whether there are jobs to be processed")
						continue
					}
				}
				log.WithError(err).Error("Error getting a job that is ready")
				continue
			}

			// TODO work happens here
			err = processDocument(job)
			if err != nil {
				log.WithError(err).WithField("job", job).Error("Error processing job")
			}

		}
	}
}
