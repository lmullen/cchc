package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var app = &App{}

func main() {
	// Check for interrupt signals and stop gracefully
	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer func() {
		signal.Stop(quit)
		cancel()
	}()
	go func() {
		select {
		case <-quit:
			log.Info("Received signal to quit the process ...")
			cancel()
		case <-ctx.Done():
		}
	}()

	err := app.Init(ctx)
	if err != nil {
		log.WithError(err).Fatal("Error initializing application")
	}
	defer app.Shutdown()

	// Perpetually run looking for items to queue for jobs, waiting in between
	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Debug("Checking for unprocessed items to add to the job queue")
			err = FindUnprocessedItems(ctx)
			if err != nil {
				log.WithError(err).Error("Error processing items with full text")
			}
			time.Sleep(15 * time.Minute)
		}
	}
}
