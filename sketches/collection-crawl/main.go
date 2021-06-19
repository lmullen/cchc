// This program crawls the LOC.gov API for a single collection and outputs basic
//  metadata about each item.
package main

import (
	"log"
)

// Hard code a collection for now
const collectionSlug = "african-american-perspectives-rare-books"
const itemsPerPage = 250

var app = &App{}

func main() {

	// Initialize the application and create a connection to the database.
	err := app.Init()
	if err != nil {
		log.Fatalln("Error initializing application: ", err)
	}
	defer app.Shutdown()

	collections, err := FetchAllCollections(app.Client)
	if err != nil {
		log.Fatalln("Error fetching all digital colletions:", err)
	}

	// Save the collections metadata to the database
	for _, c := range collections {
		err = c.Save()
		if err != nil {
			log.Println(err)
		}
	}

	// A channel to hold each page of the collection results
	collectionPages := make(chan CollectionAPIPage, 20)

	// Fetch the first page of the collection. As long as there are more pages,
	// the function will continue to fetch those too and add them to the channel.
	go fetchCollectionResult(CollectionURL(collectionSlug, 1), app.Client, collectionPages)

	// Iterate over the pages in the collection API, and the items within each page.
	// Store those results to the database.
	for r := range collectionPages {
		log.Println(r)
		for _, v := range r.Results {
			err = v.Save()
			if err != nil {
				log.Println(err)
			}
		}
	}

}
