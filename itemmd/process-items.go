package main

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// StartProcessingItems gets unfetched items
func StartProcessingItems(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// Repeat endlessly unless context is canceled
	for {
		timeout, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		unfetched, err := app.ItemsRepo.GetAllUnfetched(timeout)
		if err != nil {
			log.WithError(err).Fatal("Error getting unfetched items from database")
		}

		if len(unfetched) == 0 {
			log.Info("No unfetched items. Sleeping until checking again")
			select {
			case <-ctx.Done():
			case <-time.After(60 * time.Second):
				log.Info("Resuming checking for unfetched items")
				break
			}
		}

		for _, id := range unfetched {

			select {
			case <-ctx.Done():
				// Break out of function if context is canceled
				return
			default:
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

				// Make sure to rate limit
				app.Limiters.Items.Take()

				log.WithField("item_id", item.ID).Debug("Fetching item from loc.gov API")

				if isResourceNotItem(item.URL.String) {
					log.WithField("item_id", id).Debug("Skipping item because it is actually a resource")
					continue
				}

				err = item.Fetch(app.Client)
				if err != nil {
					log.WithError(err).WithField("item_id", id).Error("Error fetching item from API")
					// Record when the last failure happened
					app.Failures[id] = time.Now()
					log.Debug(app.Failures)
					continue
				}

				err = app.ItemsRepo.Save(ctx, item)
				if err != nil {
					log.WithError(err).WithField("item_id", id).Error("Error saving item to database")
					continue
				}
			}
		}
	}
}
