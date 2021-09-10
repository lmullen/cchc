package main

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

func processBatch(path string) {
	log.WithField("file", path).Info("Processing file")
	f, err := os.Open(path)
	if err != nil {
		log.WithField("file", path).WithError(err).Error("Error opening file")
	}

	decoder := json.NewDecoder(f)

	for decoder.More() {
		var book Book
		err = decoder.Decode(&book)
		if err != nil {
			log.Error(err)
			break
		}

		book.Year, err = year(extractDateString(book.PublicationDate))
		if err != nil {
			log.WithField("book", book).WithError(err).Warn("Error parsing date")
		}

		exists, err := book.Exists()
		if err != nil {
			log.Error(err)
			break
		}
		if !exists {
			log.WithField("book", book).Debug("Saving book to the database")
			err := book.Save()
			if err != nil {
				log.WithField("book", book).WithError(err).Error("Error saving to database")
			}
			book.Text = ""
			err = book.SaveJSON()
			if err != nil {
				log.WithField("book", book).WithError(err).Error("Error saving book JSON to database")
			}

		} else {
			log.WithField("book", book).Warn("Skipping book is already in database")
		}
	}

}
