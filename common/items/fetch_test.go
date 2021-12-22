package items

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItem_Fetch(t *testing.T) {
	item := &Item{
		ID:  "http://www.loc.gov/item/amss.hc00032b",
		URL: sql.NullString{String: "https://www.loc.gov/item/amss.hc00032b", Valid: true},
	}
	assert.False(t, item.API.Valid, "API field is not valid before fetching")
	assert.False(t, item.Fetched(), "fetched method returns false before fetching")

	err := item.Fetch(http.DefaultClient)
	assert.NoError(t, err, "fetching does not result in an error")
	assert.True(t, item.Fetched(), "fetched method returns true after fetching")

	assert.Contains(t, item.Languages, "english", "this item's language is english")

	assert.Len(t, item.Files, 7, "there are four files associated")
}
