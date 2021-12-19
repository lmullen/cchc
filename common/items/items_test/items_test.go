package items_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/lmullen/cchc/common/db"
	"github.com/lmullen/cchc/common/items"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/assert"
)

func TestItemsDB(t *testing.T) {
	t.Parallel()

	user := "gnomock"
	pass := "strong-passwords-are-the-best"
	dbname := "cchc_gnomock_test"

	p := postgres.Preset(
		postgres.WithUser(user, pass),
		postgres.WithDatabase(dbname),
	)

	container, err := gnomock.Start(p)
	assert.NoError(t, err)
	defer func() { assert.NoError(t, gnomock.Stop(container)) }()

	connstr := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		user, pass, container.Host, container.DefaultPort(), dbname)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, _ := db.Connect(ctx, connstr)
	db.Ping(ctx)
	m, _ := migrate.New("file://../../../db/migrations", connstr)
	m.Up()

	itemsRepo := items.NewItemRepo(db)

	// Create an item from ID and URL
	itemStart := &items.Item{
		ID: "http://www.loc.gov/item/mal1285100/",
		URL: sql.NullString{
			String: "https://www.loc.gov/item/mal1285100/",
			Valid:  true,
		},
	}

	// Save it to the repository
	err = itemsRepo.Save(ctx, itemStart)
	assert.NoError(t, err)

	// Retrieve the item and check that it is the same as what we started with
	itemSaved, err := itemsRepo.Get(ctx, itemStart.ID)
	itemStart.Updated = itemSaved.Updated // Set timestamp equal
	assert.True(t, reflect.DeepEqual(itemStart, itemSaved))

	// Check that the item is not fetched
	assert.False(t, itemSaved.Fetched())

	// Fetch the item from the API
	err = itemSaved.Fetch(http.DefaultClient)
	assert.NoError(t, err)

	// Check a few fields
	assert.NotEmpty(t, itemSaved.API)
	assert.Equal(t, itemSaved.Title.String, "Abraham Lincoln papers: Series 1. General Correspondence. 1833-1916: George D. Prentice, James Guthrie, and James Speed to Abraham Lincoln, Tuesday, November 05, 1861 (Telegram regarding military affairs in Kentucky)")
	assert.Equal(t, itemSaved.Year.Int32, int32(1861))

	// Save the item to the database again
	err = itemsRepo.Save(ctx, itemSaved)
	assert.NoError(t, err)

	// Check that the items are equivalent
	itemSavedAndFetched, err := itemsRepo.Get(ctx, itemSaved.ID)
	assert.NoError(t, err)

	// Check that updating the timestamp works correctly
	assert.True(t, itemSavedAndFetched.Updated.After(itemSaved.Updated))

	// PostgreSQL will change the JSONB column formatting, so don't check it
	itemSaved.API = sql.NullString{}
	itemSavedAndFetched.API = sql.NullString{}

	// Ignore the timestamp because that will by definition have changed
	itemSaved.Updated = itemSavedAndFetched.Updated

	assert.Equal(t, itemSaved, itemSavedAndFetched)
}
