package main

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// ProcessUnfetched gets unfetched items
func ProcessUnfetched(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// Repeat the cycle of fetching items endlessly until canceled
checkForUnfetched:
	for {

		// First check if we have unfetched items
		log.Info("Checking if there are any unfetched items in the database")
		check, unfetched, err := getUnfetched(ctx)
		if err != nil {
			log.WithError(err).Fatal("Error getting unfetched items from database")
			return
		}

		// If there is nothing to fetch, then wait to check again
		if !check {
			log.Info("No unfetched items that are currently checkable; will check again in a while")
			select {
			case <-ctx.Done():
				// Break out of the function if the context was canceled
				log.WithContext(ctx).Info("Work canceled: stopping processing of unfetched items")
				return
			case <-time.After(10 * time.Minute):
				log.Info("Resuming checking for unfetched items")
				continue checkForUnfetched
			}
		}

		log.WithField("unfetched", len(unfetched)).Info("Found unfetched items and starting to fetch them")
		for _, id := range unfetched {

			// Check for context cancellation before fetchign each item
			select {
			case <-ctx.Done():
				// Break out of function if context is canceled
				log.WithContext(ctx).Info("Work canceled: stopping processing of unfetched items")
				return
			default:
				// Do work on each item
				// If an item previously failed less than an hour ago skip it
				if !checkable(app.Failures, id) {
					log.WithField("item_id", id).Debug("Skipping item because it failed to fetch less than an hour ago")
					continue
				}

				// Get the item from the database and fetch it, then save to repository
				item, err := app.ItemsRepo.Get(ctx, id)
				if err != nil {
					log.WithError(err).WithField("item_id", id).Error("Error getting item to fetch from database")
					continue
				}

				if isResourceNotItem(item.URL.String) {
					log.WithField("item_id", id).Debug("Skipping item because it is actually a resource")
					continue
				}

				// Make sure to rate limit
				app.Limiters.Items.Take()
				log.WithField("item_id", item.ID).Debug("Fetching item from loc.gov API")
				err = item.Fetch(app.Client)
				if err != nil {
					log.WithError(err).WithField("item_id", id).Error("Error fetching item from API")
					// Record when the last failure happened
					app.Failures[id] = time.Now()
					continue
				} else {
					// Make sure to delete the key from the map if it is successfully fetched
					delete(app.Failures, item.ID)
				}

				timeout, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancelTimeout()
				err = app.ItemsRepo.Save(timeout, item)
				if err != nil {
					log.WithError(err).WithField("item_id", id).Error("Error saving item to database")
					continue
				}
			}
		}
	}
}
