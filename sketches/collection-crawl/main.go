// This program crawls the LOC.gov API for a single collection and outputs basic
//  metadata about each item.
package main

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

const apiBase = "https://www.loc.gov"
const itemsPerPage = 250 // TODO so far 250 seems like the best value

var app = &App{}

func main() {

	// Initialize the application and create a connection to the database.
	err := app.Init()
	if err != nil {
		log.Fatal("Error initializing application: ", err)
	}
	defer app.Shutdown()

	collections, err := FetchAllCollections(app.Client)
	if err != nil {
		log.Fatal("Error fetching all digital collections:", err)
	}

	// A channel to hold each page of the collection results
	collectionPages := make(chan CollectionAPIPage, 1000)

	// Save the collections metadata to the database, then start fetching each collection's items
	for _, c := range collections {

		// Save a collection's metadata to the database
		err = c.Save()
		if err != nil {
			log.Error(err)
		}

		// Start fetching that collection's metadata
		// TODO remove this limit which crawls only small collections
		if c.Count < 10000 {
			// Fetch the first page of the collection. As long as there are more pages,
			// the function will continue to fetch those too and add them to the channel.
			app.CollectionsWG.Add(1)
			go fetchCollectionResult(CollectionURL(c.ItemsURL, 1), c.ID, app.Client, collectionPages)
		}

	}

	// Iterate over the pages in the collection API, and the items within each page.
	// Store those results to the database. The wait group makes sure that the
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
