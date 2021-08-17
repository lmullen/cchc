// This program fetches the full metadata for LOC.gov items and writes it to the
// database. Items have previously been identified by the crawler. These items
// are then pulled off the message queue and fetched.
package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// var removeFromResponse = []string{
// 	"aka", "breadcrumbs", "browse", "categories", "content", "content_is_post",
// 	"expert_resources", "facet_trail", "facet_views", "facets", "featured_items",
// 	"form_facets", "legacy-url", "next", "next_sibling", "options",
// 	"original_formats", "pages", "partof", "previous", "previous_sibling",
// 	"research-centers", "shards", "site_type", "subjects", "timeline_1852_1880",
// 	"timeline_1881_1900", "timeline_1901_1925", "timestamp", "topics", "views",
// }

var app = &App{}

func main() {

	err := app.Init()
	if err != nil {
		log.Fatal("Error initializing application: ", err)
	}
	defer app.Shutdown()

	// Process the items from the queue
	go startProcessingItems()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
