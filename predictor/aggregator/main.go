// This program fetches the full metadata for LOC.gov items and writes it to the
// database. Items have previously been identified by the crawler. These items
// are then pulled off the message queue and fetched.
package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

var app = &App{}

func main() {

	err := app.Init()
	if err != nil {
		log.Fatal("Error initializing application: ", err)
	}
	defer app.Shutdown()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Process the items from the queue
	startProcessingDocs(context.TODO())

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit

	// log.Info("Shutdown signal received; waiting for ongoing work to finish")
	// cancel()
	// wg.Wait()

}
