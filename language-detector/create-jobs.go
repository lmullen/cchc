package main

import (
	"context"
	"sync"
	"time"

	"github.com/lmullen/cchc/common/jobs"
	log "github.com/sirupsen/logrus"
)

func createJobs(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info("Checking whether jobs need to be created")
	for {
		select {
		case <-ctx.Done():
			log.Info("Stopped creating jobs")
			return
		default:
			timeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			job, err := app.JobsRepo.CreateJobForUnqueued(timeout, queue)
			if err != nil {
				if err == jobs.ErrAllQueued {
					log.Infof("All jobs for language are queued, so waiting %s to check again", waittime)
					select {
					case <-ctx.Done():
						return
					case <-time.After(waittime):
						log.Info("Checking again whether jobs need to be created")
						continue
					}
				}
				log.WithError(err).Error("Error creating job")
				continue
			}
			log.WithField("job", job).Debug("Created job")
		}
	}
}
