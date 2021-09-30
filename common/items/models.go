package items

import (
	"context"
	"database/sql"
)

// Item is a representation of an item in the LOC digital collections, along
// with its resources and files.
type Item struct {
	ID        string
	URL       sql.NullString
	Title     sql.NullString
	Year      sql.NullInt32
	Date      sql.NullString
	Subjects  []string
	Timestamp sql.NullInt64
	Resources []ItemResource
	Files     []ItemFile
	API       sql.NullString // The entire API response stored as JSONB
}

// Repository is an interface describing a data store for items.
type Repository interface {
	Get(ctx context.Context, ID string) (*Item, error)
	// Save(item *Item) error
}

// ItemResource is a resource attached to an item.
type ItemResource struct {
	ItemID       string
	ResourceSeq  int
	FullTextFile sql.NullString
	DJVUTextFile sql.NullString
	Image        sql.NullString
	PDF          sql.NullString
	URL          sql.NullString
	Caption      sql.NullString
}

// ItemFile is a file contained within an item. Unlike the LOC.gov API, this
// model does not make a firm distinction between a file and format.
type ItemFile struct {
	ItemID          string
	ResourceSeq     int
	FileSeq         int
	FormatSeq       int
	Mimetype        sql.NullString
	FullText        sql.NullString
	FullTextService sql.NullString
	WordCoordinates sql.NullString
	URL             sql.NullString
	Info            sql.NullString
	Use             sql.NullString
}
