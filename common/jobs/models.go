package jobs

import (
	"database/sql"

	"github.com/google/uuid"
)

// FullText is a record of how jobs with the full text of items were passed on to
// some particular destination. It can handle recording full text jobs at either
// the item (i.e., resource) or file level, and it can handle multiple destinations.
type FullText struct {
	ID          uuid.UUID
	ItemID      string
	Destination string
	Started     sql.NullTime
	Finished    sql.NullTime
	Status      string
}
