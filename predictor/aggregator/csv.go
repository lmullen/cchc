package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/lmullen/cchc/common/results"

	"github.com/lmullen/cchc/common/messages"
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
func processPredictionsCSV(ctx context.Context, path string) error {
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
		// log.Debug(p)
		jobID, err := uuid.Parse(p[0])
		if err != nil {
			return err
		}
		prob, err := strconv.ParseFloat(p[4], 64)
		if err != nil {
			return err
		}
		q := results.NewQuotation(jobID, p[1], p[2], p[3], prob)
		err = app.ResultsRepo.Save(ctx, q)
		if err != nil {
			return err
		}
	}

	return nil
}
