package main

import (
	"context"
	"fmt"
	"time"
)

// Get unfetched returns a boolean saying whether to do work or not, a slice
// of unfetched items, and any errors.
func getUnfetched(ctx context.Context) (bool, []string, error) {
	timeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	unfetched, err := app.ItemsRepo.GetAllUnfetched(timeout)

	// If there is a problem checking for unfetched items
	if err != nil {
		return false, nil, fmt.Errorf("Error getting unfetched items from database: %w", err)
	}

	// If there are no unfetched items
	if len(unfetched) == 0 {
		return false, unfetched, nil
	}

	// See how many uncheckable items we have in the list of unfetched. If we have
	// more uncheckable than are to be fetched, all the ones to be fetched must be
	// previous failures. So skip for now.
	uncheckable := 0
	for id := range app.Failures {
		if !checkable(app.Failures, id) {
			uncheckable++
		}
	}
	if uncheckable >= len(unfetched) {
		return false, unfetched, nil
	}

	// Return that we should check items and return the slice of IDs
	return true, unfetched, nil

}
