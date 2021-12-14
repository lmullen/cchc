package main

import (
	"context"
	"sort"

	"github.com/mpvl/unique"
)

// FindUnprocessedItems takes a query as input. That query should return a list
// of item IDs which should be enqueued but which do not have a job in the
// table. Those items will then be handled in a separate function.
func FindUnprocessedItems(ctx context.Context, query string) ([]string, error) {
	var items []string

	rows, err := app.DB.Query(ctx, query)
	defer rows.Close()

	for rows.Next() {
		var item string
		err := rows.Scan(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	sort.Strings(items)
	unique.Strings(&items)

	return items, nil
}
