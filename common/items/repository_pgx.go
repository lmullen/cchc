package items

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Repo is a data store using PostgreSQL with the pgx native interface.
type Repo struct {
	db *pgxpool.Pool
}

// NewItemRepo returns an item repo using PostgreSQL with the pgx native interface.
func NewItemRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

// Get fetches an item from the database by its ID.
func (r *Repo) Get(ctx context.Context, ID string) (*Item, error) {
	item := Item{}
	itemQuery := `
	SELECT id, url, title, year, date, subjects, languages, api
	FROM items 
	WHERE id = $1;
	`

	resourcesQuery := `
	SELECT item_id, resource_seq, fulltext_file, djvu_text_file, image, pdf, url, caption
	FROM resources 
	WHERE item_id = $1;
	`

	filesQuery := `
	SELECT item_id, resource_seq, file_seq, format_seq, mimetype, fulltext, fulltext_service, word_coordinates, url, info, use
	FROM files
	WHERE item_id = $1;
	`

	err := r.db.QueryRow(ctx, itemQuery, ID).
		Scan(&item.ID, &item.URL, &item.Title, &item.Year, &item.Date, &item.Subjects, &item.Languages, &item.API)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, resourcesQuery, ID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		res := ItemResource{}
		err = rows.Scan(&res.ItemID, &res.ResourceSeq, &res.FullTextFile, &res.DJVUTextFile, &res.Image, &res.PDF, &res.URL, &res.Caption)
		if err != nil {
			return nil, err
		}
		item.Resources = append(item.Resources, res)
	}
	rows.Close()

	rows, err = r.db.Query(ctx, filesQuery, ID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		res := ItemFile{}
		err = rows.Scan(&res.ItemID, &res.ResourceSeq, &res.FileSeq, &res.FormatSeq, &res.Mimetype, &res.FullText, &res.FullTextService, &res.WordCoordinates, &res.URL, &res.Info, &res.Use)
		if err != nil {
			return nil, err
		}
		item.Files = append(item.Files, res)
	}
	rows.Close()

	return &item, nil
}

// Save serializes an item to the database, either creating it in the database
// or updating the fields.
func (r *Repo) Save(ctx context.Context, item *Item) error {
	itemQuery := `
	INSERT INTO items (id, url, title, year, date, subjects, api, languages)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (id) DO UPDATE
	SET
	  url              = $2,
		title            = $3,
		year             = $4,
		date             = $5,
		subjects         = $6,
		api              = $7,
		languages        = $8;
	`

	resourceQuery := `
	INSERT INTO resources (item_id, resource_seq, fulltext_file, djvu_text_file,
		image, pdf, url, caption)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	fileQuery := `
	INSERT INTO files (item_id, resource_seq, file_seq, format_seq,
	                   mimetype, fulltext, fulltext_service, word_coordinates,
										 url, info, use)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	// Use a transaction since we are writing to three tables
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Error creating transaction in database: %w", err)
	}

	_, err = tx.Exec(ctx, itemQuery, item.ID, item.URL, item.Title, item.Year,
		item.Date, item.Subjects, item.API, item.Languages)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("Error saving item %s to database: %w", item, err)
	}

	for _, r := range item.Resources {
		_, err = tx.Exec(ctx, resourceQuery, r.ItemID, r.ResourceSeq, r.FullTextFile,
			r.DJVUTextFile, r.Image, r.PDF, r.URL, r.Caption)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("Error saving item %s to database: %w", item, err)
		}
	}

	for _, f := range item.Files {
		_, err = tx.Exec(ctx, fileQuery, f.ItemID, f.ResourceSeq, f.FileSeq, f.FormatSeq,
			f.Mimetype, f.FullText, f.FullTextService, f.WordCoordinates, f.URL,
			f.Info, f.Use)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("Error saving item %s to database: %w", item, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Error saving item %s to database: %w", item, err)
	}

	return nil
}
