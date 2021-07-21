// This program crawls the LOC.gov API and retrieves information about its
// digital collections. It proceeds in this way.
//   1. It fetches all digital collections (with some filtering)
//   2. It fetches the items in those digital collections via the search (again
//      with some filtering).
//   3. It then fetches metadata about those items from the API.
//
// Everything gets stored in a database.
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Configuration options that aren't worth exposing as environment variables
const (
	apiBase         = "https://www.loc.gov"
	apiItemsPerPage = 1000
	apiTimeout      = 60 // The timeout limit for API requests in seconds
)

var removeFromResponse = []string{
	"aka", "breadcrumbs", "browse", "categories", "content", "content_is_post",
	"expert_resources", "facet_trail", "facet_views", "facets", "featured_items",
	"form_facets", "legacy-url", "next", "next_sibling", "options",
	"original_formats", "pages", "partof", "previous", "previous_sibling",
	"research-centers", "shards", "site_type", "subjects", "timeline_1852_1880",
	"timeline_1881_1900", "timeline_1901_1925", "timestamp", "topics", "views",
}

var app = &App{}

func main() {

	// Initialize the application and create a connection to the database.
	err := app.Init()
	if err != nil {
		log.Fatal("Error initializing application: ", err)
	}
	defer app.Shutdown()

	// Set the logging level
	switch app.Config.loglevel {
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	}

	// A channel to hold each page of the collection results
	collectionPages := make(chan CollectionAPIPage, 1000)

	// In a goroutine, fetch all the digital collections periodically. This is the
	// entry point: all the collections will be detected, and then all the items
	// in those collections.
	go func() {

		for { // This will happen forever in a goroutine until the program is quit
			collections, err := FetchAllCollections()
			if err != nil {
				log.Error("Error fetching all digital collections:", err)
				time.Sleep(1 * time.Hour) // Don't wait forever, but don't try again right away
				break                     // Start over trying to fetch all collections
			}

			// Save the metadata for each collection to the database, then start fetching each
			// collection's items
			for _, c := range collections {

				// Save a collection's metadata to the database
				err = c.Save()
				if err != nil {
					log.WithField("collection", c).Error("Error saving collection to database:", err)
				}

				// Start fetching the items in that collection.
				// TODO remove this limit which crawls only small collections
				if c.Count < 50 {
					// Fetch the first page of the collection. As long as there are more pages,
					// the function will continue to fetch those too and add them to the channel.
					go c.FetchCollectionItems(1, collectionPages)
				}

			}

			// Goroutines have been started for fetching each collections items. We
			// want to wait a decent interval, and then start the crawl over again
			// from the beginning.
			// TODO: Set this to some reasonable default, like a week
			time.Sleep(2 * time.Minute)
			// Now the loop starts over again by fetching all the digital collections
		}

	}()

	// Iterate over the pages in the collection API, and the items within each
	// page. Store those results to the database. This means we know that an item
	// exists, and also which collection it is a part of. But we will fetch that
	// item from its item page separately.
	go func() {
		// This will effectively iterate forever, because the channel will not be
		// closed. But it will also not do work unless there are pages in the
		// channel.
		for r := range collectionPages {
			// Start a new goroutine to deal with each page
			go func(r CollectionAPIPage) {
				for _, item := range r.Results {
					item.CollectionID = r.CollectionID
					err = item.Save()
					if err != nil {
						log.WithFields(log.Fields{
							"item_id": item.ID,
							"error":   err,
						}).Error("Error saving item")
					}

					item := item.ToItem()
					fetched, err := item.Fetched()
					if err != nil {
						log.Error("Error checking if item has been fetched: ", err)
					}
					if fetched {
						// Don't put the message in the queue if we've already fetched it
						return
					}
					err = item.EnqueueMetadata()
					if err != nil {
						log.WithFields(log.Fields{
							"item_id": item.ID,
							"error":   err,
						}).Error("Error putting item in queue for metadata processing")
					}
				}
			}(r)
		}
	}()

	// Process the items from the queue
	go func() {
		for msg := range app.ItemMetadataQ.Consumer {
			// Give each item its own goroutine
			go func(msg amqp.Delivery) {
				err = ProcessItemMetadata(msg)
				if err != nil {
					log.Error("Error processing item from queue: ", err)
				}
			}(msg)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	close(collectionPages)

}
