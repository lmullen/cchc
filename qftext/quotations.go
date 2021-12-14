package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

func enqueueForQuotations(ctx context.Context) {
	// A query that gets American history collections
	query := `
	SELECT i.id
  FROM items i
	LEFT JOIN items_in_collections ic
	ON i.id = ic.item_id
	LEFT JOIN collections c
	ON ic.collection_id = c.id
  WHERE i.api IS NOT NULL 
		AND 'american history' = ANY(c.topics)
		AND NOT (EXISTS ( SELECT
          FROM jobs.fulltext
          WHERE fulltext.item_id = i.id AND queue = 'fulltext-quotations'));
	`

	// Do this perpetually, sleeping between jobs
	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Info("Putting any American history items with full text into the queue")

			items, err := FindUnprocessedItems(ctx, query)
			if err != nil {
				log.WithError(err).Error("Problem finding unprocessed items for biblical quotations")
				return
			}

			for _, item := range items {
				select {
				case <-ctx.Done():
					return
				default:
					err = EnqueueItemFulltext(ctx, item, app.MR.quotations)
					if err != nil {
						log.WithError(err).Error("Problem sending full text of item")
					}
				}

			}
		}

		time.Sleep(10 * time.Minute)
	}

}
