package items

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// ItemRepo is a data store using PostgreSQL with the pgx native interface.
type ItemRepo struct {
	db *pgxpool.Pool
}

// NewItemRepo returns an item repo using PostgreSQL with the pgx native interface.
func NewItemRepo(db *pgxpool.Pool) *ItemRepo {
	return &ItemRepo{
		db: db,
	}
}

// Get fetches an item from the database by its ID.
func (r *ItemRepo) Get(ctx context.Context, ID string) (*Item, error) {
	item := Item{}
	itemQuery := `
	SELECT id, url, title, year, date, subjects, timestamp, api
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
		Scan(&item.ID, &item.URL, &item.Title, &item.Year, &item.Date, &item.Subjects, &item.Timestamp, &item.API)
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
// func (r *ItemRepo) Save(ctx context.Context, item *Item) error {
// 	return nil
// }
