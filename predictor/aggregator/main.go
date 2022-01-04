// This program gets a batch
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

const queue = "quotations"
const waittime = 15 * time.Minute
const jobtimeout = 30 * time.Minute

var app = &App{}

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

	err := app.Init()
	if err != nil {
		log.Fatal("Error initializing application: ", err)
	}
	defer app.Shutdown()

	wg := &sync.WaitGroup{}

	// Create jobs for items
	wg.Add(1)
	go createJobs(ctx, wg)

	// Process the items from the queue
	wg.Add(1)
	// Sleep a bit to give time for the jobs to be created before processing
	time.Sleep(15 * time.Second)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Info("Stopped processing batches")
				return
			default:
				processBatchOfDocs(ctx)
			}
		}
	}()

	wg.Wait()

}
