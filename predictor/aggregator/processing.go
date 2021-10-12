package main

import (
	"encoding/json"
	"os"

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
		fDocs, err := writeDocsCSV(docs)
		if err != nil {
			log.WithError(err).Error("Error writing CSV to send to prediction model")
		}
		log.Debug(fDocs)

		// TODO: Call the prediction model

		// Get the predictions back from a temporary file and write them to the database
		err = processPredictionsCSV(fDocs)
		if err != nil {
			log.WithError(err).Error("Error getting results from prediction model")
		}
		// Acknowledge the messages in the batch
		for _, m := range msgs {
			m.Ack(false)
		}

		// Clean up the temporary files
		err = os.Remove(fDocs)
		if err != nil {
			log.WithError(err).Warn("Problem removing the temporary files")
		}

	}
}
