package jobs

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lmullen/cchc/common/items"
)

func (job FulltextPredict) String() string {
	return fmt.Sprintf("Job %s for item %s ", job.ID, job.ItemID)
}

// Create fills out the details of the struct
func (job *FulltextPredict) Create(itemID string) {
	// Check to see what kind of full text we have available
	job.ID = uuid.New()
	job.ItemID = itemID
}

func (job *FulltextPredict) FulltextFromFile(file items.ItemFile, source string) {
	job.HasFTMethod = true
	job.ResourceSeq.Scan(file.ResourceSeq)
	job.FileSeq.Scan(file.FileSeq)
	job.FormatSeq.Scan(file.FormatSeq)
	job.Level.Scan("file")
	job.Source.Scan(source)
}

// Start records that a job has a full text method and records the start time.
func (job *FulltextPredict) Start() {
	job.HasFTMethod = true
	job.Started.Scan(time.Now())
}

// Finish records the finishing time for a job.
func (job *FulltextPredict) Finish() error {
	err := job.Finished.Scan(time.Now())
	if err != nil {
		return err
	}
	return nil
}
