// This program gets a batch
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

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
	wg.Add(1)

	// Process the items from the queue
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				processBatchOfDocs(ctx)
			}
		}
	}()

	wg.Wait()

}
