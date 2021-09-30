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

// Create fills out the details of the struct. If start is true, a started time
// will be recorded; otherwise it will be null.
func (job *FulltextPredict) Create(itemID string, start bool) {
	// Check to see what kind of full text we have available
	job.ID = uuid.New()
	job.ItemID = itemID
	if start {
		job.Started.Scan(time.Now())
	}
}

// Finish records the finishing time for a job.
func (job *FulltextPredict) Finish() error {
	err := job.Finished.Scan(time.Now())
	if err != nil {
		return err
	}
	return nil
}

// PlainTextFullText gets the plain text from a file level fulltext field if the
// item has a plaintext mimetype. It returns the full text that will be used.
func (job *FulltextPredict) PlainTextFullText(file items.ItemFile) string {
	job.HasFTMethod = true
	job.ResourceSeq.Scan(file.ResourceSeq)
	job.FileSeq.Scan(file.FileSeq)
	job.FormatSeq.Scan(file.FormatSeq)
	job.Level.Scan("file")
	job.Source.Scan("Mimetype: text/plain; Source: fulltext.")
	return file.FullText.String
}

// XMLFullText gets the plain text from a file level fulltext field if the item
// has an XML mimetype. It returns the full text that will be used.
func (job *FulltextPredict) XMLFullText(file items.ItemFile) string {
	job.HasFTMethod = true
	job.ResourceSeq.Scan(file.ResourceSeq)
	job.FileSeq.Scan(file.FileSeq)
	job.FormatSeq.Scan(file.FormatSeq)
	job.Level.Scan("file")
	job.Source.Scan("Mimetype: text/xml; Source: fulltext.")
	return file.FullText.String
}
