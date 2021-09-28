package jobs

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (j JobFulltextPredict) String() string {
	return fmt.Sprintf("Job %s for item %s ", j.ID, j.ItemID)
}

// Create fills out the details of the struct
func (j *JobFulltextPredict) Create(itemID string) error {
	// Check to see what kind of full text we have available
	j.ID = uuid.New()
	j.ItemID = itemID
	err := j.Started.Scan(time.Now())
	if err != nil {
		return err
	}
	return nil
}
