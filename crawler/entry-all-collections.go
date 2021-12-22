package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// StartFetchingCollections will fetch the digital collections, pass the pages
// into a channel, and then start over again at the proper crawl interval. It
// provides an entry point into the LOC.gov API. It is intended to be run in its
// own goroutine.
func StartFetchingCollections(cp chan CollectionAPIPage) {
	for { // This will happen forever until the program is quit
		log.Info("Starting a crawl of all collections")
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
			// Fetch the first page of the collection. As long as there are more pages,
			// the function will continue to fetch those too and add them to the channel.
			go c.FetchCollectionItems(1, cp)

		}

		// Goroutines have been started for fetching each collections items. We
		// want to wait a decent interval, and then start the crawl over again
		// from the beginning.
		log.Infof("Waiting to start a crawl of all collections for %s", crawlInterval)
		time.Sleep(crawlInterval)
		// Now the loop starts over again by fetching all the digital collections
	}

}

// StartProcessingCollections uses a channel of pages from the collections and
// processes each page, saving the items to the database and passing them to the
// message queue.
func StartProcessingCollections(cp chan CollectionAPIPage) {

	// This will effectively iterate forever, because the channel will not be
	// closed. But it will also not do work unless there are pages in the
	// channel.
	for r := range cp {
		// Start a new goroutine to deal with each page
		go func(r CollectionAPIPage) {
			for _, item := range r.Results {
				item.CollectionID = r.CollectionID
				err := item.Save()
				if err != nil {
					log.WithFields(log.Fields{
						"item_id": item.ID,
						"error":   err,
					}).Error("Error saving item")
				}
			}
		}(r)
	}
}
