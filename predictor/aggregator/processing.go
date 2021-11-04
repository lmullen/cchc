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

func processBatchOfDocs(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		var msgs []*amqp.Delivery
		var docs []*messages.FullTextPredict
		numInBatch := 100

		log.Debugf("Processing a batch of %v documents", numInBatch)

		// Read a batch of full text
		for i := 0; i < numInBatch; i++ {
			var doc messages.FullTextPredict
			msg := <-app.MsgRepo.Consume()
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

		ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx,
			"Rscript", "/predictor/id-quotations.R",
			"--bible", "bible-payload.rda",
			"--model", "prediction-payload.rda",
			"--verbose", "2",
			"--out", predictionsFile.Name(),
			"--potential",
			docsFile,
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.WithError(err).WithField("R-output", string(output)).Error("Problem running prediction model in R")
		}

		// Get the predictions back from a temporary file and write them to the database
		err = processPredictionsCSV(ctx, predictionsFile.Name())
		if err != nil {
			log.WithError(err).Error("Error getting results from prediction model")
		}

		// Acknowledge the messages in the batch
		for _, m := range msgs {
			m.Ack(false)
		}

		//
		for _, d := range docs {
			job, err := app.JobsRepo.Get(ctx, d.JobID)
			if err != nil {
				log.WithError(err).WithField("job-id", d.JobID).Warning("Problem getting job from database")
				break
			}
			job.Finish()
			err = app.JobsRepo.Save(ctx, job)
			if err != nil {
				log.WithError(err).WithField("job-id", d.JobID).Warning("Problem updating finished job in database")
			}
		}

		// Clean up the temporary files
		err = os.Remove(predictionsFile.Name())
		if err != nil {
			log.WithError(err).Warn("Problem removing the temporary files")
		}
		err = os.Remove(docsFile)
		if err != nil {
			log.WithError(err).Warn("Problem removing the temporary files")
		}

		log.Debugf("Finished processing a batch of %v documents", numInBatch)

	}
}
