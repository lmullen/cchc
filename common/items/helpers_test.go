package items

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_year(t *testing.T) {
	assert.Equal(t, sql.NullInt32{Int32: 1980, Valid: true}, year("1980-09-09"))
	assert.False(t, year("19??-09-09").Valid)
}
