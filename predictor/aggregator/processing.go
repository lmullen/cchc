package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"time"

	"github.com/lmullen/cchc/common/messages"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

func startProcessingDocs(ctx context.Context) {
	// TODO: Eventually this should run perpetually
	select {
	case <-ctx.Done():
		return
	default:
		var msgs []*amqp.Delivery
		var doc messages.FullTextPredict
		var docs []*messages.FullTextPredict
		numInBatch := 10

		// Read a batch of full text
		for i := 0; i < numInBatch; i++ {
			msg := <-app.DocumentsQ.Consumer
			err := json.Unmarshal(msg.Body, &doc)
			if err != nil {
				log.WithError(err).Error("Error processing doc")
				msg.Reject(false)
			}
			docs = append(docs, &doc)
			msgs = append(msgs, &msg)
		}

		// Write the full text to a temporary CSV
		docsFile, err := writeDocsCSV(docs)
		if err != nil {
			log.WithError(err).Error("Error writing CSV to send to prediction model")
		}

		// Create a temp file for output.
		predictionsFile, err := os.CreateTemp("", "prediction-*.csv")
		if err != nil {
			log.WithError(err).Error("Error creating temporary file for predictions")
		}
		predictionsFile.Close()

		// TODO: Call the prediction model
		ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx,
			"Rscript", "/predictor/id-quotations.R",
			"--bible", "bible-payload.rda",
			"--model", "prediction-payload.rda",
			"--verbose", "0",
			"--out", predictionsFile.Name(),
			docsFile,
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.WithError(err).WithField("R-output", string(output)).Error("Problem running prediction model in R")
		}

		// Get the predictions back from a temporary file and write them to the database
		err = processPredictionsCSV(predictionsFile.Name())
		if err != nil {
			log.WithError(err).Error("Error getting results from prediction model")
		}
		// Acknowledge the messages in the batch
		for _, m := range msgs {
			m.Ack(false)
		}

		// Clean up the temporary files
		// err = os.Remove(predictionsFile.Name())
		// if err != nil {
		// 	log.WithError(err).Warn("Problem removing the temporary files")
		// }
		// err = os.Remove(docsFile)
		// if err != nil {
		// 	log.WithError(err).Warn("Problem removing the temporary files")
		// }

	}
}
