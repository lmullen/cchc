package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

var app = &App{}

func main() {
	err := app.Init()
	if err != nil {
		log.WithError(err).Fatal("Error initializing application")
	}
	defer app.Shutdown()

	err = ProcessUnqueued(context.TODO())
	if err != nil {
		log.WithError(err).Error("Error processing unenqueued full text items")
	}

}
