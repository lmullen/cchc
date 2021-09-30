package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var app = &App{}

func main() {
	err := app.Init()
	if err != nil {
		log.WithError(err).Fatal("Error initializing application")
	}
	defer app.Shutdown()

	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
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

	err = FindUnprocessedItems(ctx)
	if err != nil {
		log.WithError(err).Fatal("Error processing items with full text")
	}

}
