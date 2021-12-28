package jobs

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (job FullText) String() string {
	return fmt.Sprintf("[Job %s. Item: %s. Status: %s.]", job.ID, job.ItemID, job.Status)
}

// NewFullText creates a new, unstarted job for full text. You must specify
// which item this is associated with, what the destination is (i.e., the job
// to be done), and whether or not it has full text associated with it. (If it
// does not have full text, presumably you will skip the job.)
func NewFullText(itemID string, destination string) *FullText {
	return &FullText{
		ID:          uuid.New(),
		ItemID:      itemID,
		Destination: destination,
		Status:      "ready",
	}
}

// Start adds the current time to the started field and changes the job status.
func (job *FullText) Start() {
	job.Started.Scan(time.Now())
	job.Status = "running"
}

// Finish adds the current time to the finished field and changes the job status.
func (job *FullText) Finish() {
	job.Finished.Scan(time.Now())
	job.Status = "finished"
}

// Skip adds the current time to the finished field and changes the job status.
func (job *FullText) Skip() {
	job.Finished.Scan(time.Now())
	job.Status = "skipped"
}

// Fail adds the current time to the finished field and changes the job status.
func (job *FullText) Fail() {
	job.Finished.Scan(time.Now())
	job.Status = "failed"
}
