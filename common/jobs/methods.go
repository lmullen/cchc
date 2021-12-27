package jobs

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (job FullText) String() string {
	var status string
	switch {
	case job.Started.Valid && job.Finished.Valid:
		status = "finished"
	case job.Started.Valid && !job.Finished.Valid:
		status = "started"
	case !job.Started.Valid && !job.Finished.Valid:
		status = "unstarted"
	}
	return fmt.Sprintf("[Job %s. Item: %s. Status: %s.]", job.ID, job.ItemID, status)
}

// NewFullText creates a new, unstarted job for full text. You must specify
// which item this is associated with, what the destination is (i.e., the job
// to be done), and whether or not it has full text associated with it. (If it
// does not have full text, presumably you will skip the job.)
func NewFullText(itemID string, destination string, hasFT bool) *FullText {
	return &FullText{
		ID:          uuid.New(),
		ItemID:      itemID,
		HasFTMethod: hasFT,
		Destination: destination,
	}
}

// Start adds the current time to the started field.
func (job *FullText) Start() {
	job.Started.Scan(time.Now())
}

// Finish adds the current time to the finished field.
func (job *FullText) Finish() {
	job.Finished.Scan(time.Now())
}
