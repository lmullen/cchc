// This program crawls the LOC.gov API and identifies the existence of digitized items
//
// It proceeds in this way.
//
// 1. It fetches all digital collections (with some filtering)
// 2. It fetches the items in those digital collections (again
//    with some filtering for full text items).
//
// Everything gets stored in a database.
//
// The task of fetching the full item metadata is handled by itemmd.
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// Configuration options that aren't worth exposing as environment variables
const (
	apiBase         = "https://www.loc.gov"
	apiItemsPerPage = 1000
	apiTimeout      = 60 // The timeout limit for API requests in seconds
)

// How long to wait between starting new crawls at the entry points
var crawlInterval time.Duration = 2 * 24 * time.Hour

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

	// A channel to hold each page of the collection results
	collectionPages := make(chan CollectionAPIPage, 1000)

	// In a goroutine, fetch all the digital collections periodically. This is the
	// entry point: all the collections will be detected, and then all the items
	// in those collections.
	go StartFetchingCollections(collectionPages)

	// Iterate over the pages in the collection API, and the items within each
	// page. Store those results to the database. This means we know that an item
	// exists, and also which collection it is a part of. But we will fetch that
	// item from its item page separately.
	go StartProcessingCollections(collectionPages)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	close(collectionPages)

}
