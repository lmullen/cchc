package main

import (
	"fmt"
)

// ItemResult is an item returned from a LOC.gov collection results page. There
// are many more fields that are returned in a collections result page, but we
// are going to get that data directly from the item page instead.
type ItemResult struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	CollectionID string // This is the foreign key to the collection, not from API
}

// Use the title as a string representation of an item.
func (item ItemResult) String() string {
	return item.ID
}

// Save serializes an item to the database.
func (item ItemResult) Save() error {
	itemQuery := `
	INSERT INTO items(id, url) 
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING;
	`

	relationQuery := `
	INSERT INTO items_in_collections(item_id, collection_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING;
	`
	itemStmt, err := app.DB.Prepare(itemQuery)
	if err != nil {
		return fmt.Errorf("Error preparing item save query: %w", err)
	}

	relationStmt, err := app.DB.Prepare(relationQuery)
	if err != nil {
		return fmt.Errorf("Error preparing item/collection query: %w", err)
	}

	// Use a transaction since we are writing to two tables
	tx, err := app.DB.Begin()
	if err != nil {
		return fmt.Errorf("Error creating transaction in database: %w", err)
	}

	_, err = tx.Stmt(itemStmt).Exec(item.ID, item.URL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error saving item to database: %w", err)
	}

	_, err = tx.Stmt(relationStmt).Exec(item.ID, item.CollectionID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error saving item/collection relation to database: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Error committing item to database: %w", err)
	}

	return nil

}

// ToItem converts an ItemResult to an Item.
func (item ItemResult) ToItem() Item {
	return Item{
		ID:  item.ID,
		URL: item.URL,
	}
}
