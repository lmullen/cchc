package jobs

import (
	"database/sql"
)

// JobFulltextPredict is a record of how full text was passed on to the
// prediction model for each item.
type JobFulltextPredict struct {
	ID          int64
	ItemID      string
	Level       string
	Source      string
	HasFTMethod bool
	started     sql.NullTime
	finished    sql.NullTime
}

func (j JobFulltextPredict) String() string {
	return "Job for item: " + j.ItemID
}

func (j *JobFulltextPredict) Create() error {
	// Check to see what kind of full text we have available
	return nil
}
