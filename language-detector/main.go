package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var app App

const queue = "language"
const waittime = 1 * time.Minute

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Clean up function that will be called at program end no matter what
	defer func() {
		signal.Stop(quit)
		cancel()
	}()
	// Listen for shutdown signals in a go-routine and cancel context then
	go func() {
		select {
		case <-quit:
			log.Info("Shutdown signal received; quitting language detector")
			cancel()
		case <-ctx.Done():
		}
	}()

	err := app.Init(ctx)
	if err != nil {
		log.Fatal("Error initializing application: ", err)
	}
	defer app.Shutdown()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go createJobs(ctx, wg)
	// go func() {
	// 	for i := 0; i < 100; i++ {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		default:
	// 			job, err := app.JobsRepo.CreateJobForUnqueued(ctx, queue)
	// 			if err != nil {
	// 				if err == jobs.ErrAllQueued {
	// 					log.Info("All jobs for language are queued. Waiting fifteen minutes.")
	// 					select {
	// 					case <-time.After(1 * time.Minute):
	// 						continue
	// 					case <-ctx.Done():
	// 						return
	// 					}
	// 				}
	// 				log.WithError(err).Error("Error creating job")
	// 				continue
	// 			}
	// 			log.WithField("job", job).Debug("Created job")
	// 		}
	// 	}
	// }()

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		default:
	// 			job, err := app.JobsRepo.GetReadyJob(ctx, queue)
	// 			if err != nil {
	// 				if err == jobs.ErrNoJobs {
	// 					log.Info("No ready jobs. Waiting fifteen minutes.")
	// 					select {
	// 					case <-time.After(1 * time.Minute):
	// 						continue
	// 					case <-ctx.Done():
	// 						return
	// 					}
	// 				}
	// 				log.WithError(err).Error("Error getting ready job")
	// 				continue
	// 			}
	// 			log.WithField("job", job).Debug("Got a ready job")
	// 			time.Sleep(2 * time.Second)
	// 			job.Skip()
	// 			err = app.JobsRepo.SaveFullText(ctx, job)
	// 			if err != nil {
	// 				log.WithError(err).Error("Error saving job after skipping it")
	// 			}
	// 			log.WithField("job", job).Debug("Skipped the job")
	// 		}
	// 	}
	// }()

	wg.Wait()

}
