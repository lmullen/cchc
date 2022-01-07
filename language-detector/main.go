// This program identifies the language of each sentence in full-text items.
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
const waittime = 15 * time.Minute
const jobtimeout = 120 * time.Second

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

	// Sleep a bit to give time for the jobs to be created before processing
	time.Sleep(15 * time.Second)
	wg.Add(1)
	go processJobs(ctx, wg)

	wg.Wait()

}
