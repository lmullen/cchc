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
	"sync"

	log "github.com/sirupsen/logrus"
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

	switch app.Config.loglevel {
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	}

	collections, err := FetchAllCollections()
	if err != nil {
		log.Error("Error fetching all digital collections:", err)
	}

	// A channel to hold each page of the collection results
	collectionPages := make(chan CollectionAPIPage, 1000)

	// Save the collections metadata to the database, then start fetching each
	// collection's items
	for _, c := range collections {

		// Save a collection's metadata to the database
		err = c.Save()
		if err != nil {
			log.WithField("collection", c).Error("Error saving collection to database: err")
		}

		// Start fetching that collection's metadata
		// TODO remove this limit which crawls only small collections
		if c.Count < 10000 {
			// Fetch the first page of the collection. As long as there are more pages,
			// the function will continue to fetch those too and add them to the channel.
			app.CollectionsWG.Add(1)
			go c.FetchCollectionItems(1, collectionPages)
		}

	}

	// Iterate over the pages in the collection API, and the items within each
	// page. Store those results to the database. This means we know that an item
	// exists, and also which collection it is a part of. But we will fetch that
	// item from its item page separately. The wait group makes sure that the
	// program doesn't end before the data is written to the database.
	itemsWG := sync.WaitGroup{}
	itemsWG.Add(1)
	go func() {
		for r := range collectionPages {
			for _, item := range r.Results {
				item.CollectionID = r.CollectionID
				err = item.Save()
				if err != nil {
					log.WithFields(log.Fields{
						"item_id": item.ID,
						"error":   err,
					}).Error("Error saving item")
				}
			}
		}
		itemsWG.Done()
	}()

	app.CollectionsWG.Wait()
	close(collectionPages) // Make sure the program quits
	itemsWG.Wait()

}
