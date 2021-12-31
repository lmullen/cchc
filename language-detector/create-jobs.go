package main

import (
	"context"
	"strings"
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
				// Checking the error by value does not work
				// if pgErr, ok := err.(*pgconn.PgError); !ok || pgErr.Code != "23505" {
				// 	log.WithError(err).Debug("Attempt to create duplicate job failed")
				// }
				// So just check the bloody thing by matching strings
				if strings.Contains(err.Error(), "SQLSTATE 23505") {
					// TODO This is an error, but not at this problem of the program. So
					// only log it if we really want to know all the dirty details. It doesn't
					// actually cause a problem and needs to be fixed with the queuing locks
					// in the jobs package.
					log.WithError(err).Trace("Attempt to create duplicate job failed")
				} else {
					log.WithError(err).Error("Error creating job")
				}
				continue
			}
			log.WithField("job", job).Debug("Created job")
		}
	}
}
