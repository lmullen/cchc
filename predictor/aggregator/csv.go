package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/lmullen/cchc/common/messages"
	log "github.com/sirupsen/logrus"
)

// Write out a CSV with the full text for the prediction model
func writeDocsCSV(docs []*messages.FullTextPredict) (string, error) {
	f, err := os.CreateTemp("", "fulltext-*.csv")
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, doc := range docs {
		err := w.Write(doc.CSVRow())
		if err != nil {
			return "", fmt.Errorf("Error writing temporary CSV: %w", err)
		}
	}

	return f.Name(), nil
}

// Read in results of the prediction model and write them to the database.
func processPredictionsCSV(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)

	predictions, err := r.ReadAll()
	if err != nil {
		return err
	}

	for _, p := range predictions {
		// TODO: This is where to do the work for each prediction
		log.Debug(p)
	}

	return nil
}
