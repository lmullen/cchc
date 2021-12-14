package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lmullen/cchc/common/db"
	"github.com/lmullen/cchc/common/items"
	"github.com/lmullen/cchc/common/jobs"
	"github.com/lmullen/cchc/common/messages"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
)

// App holds resources and config
type App struct {
	DB *pgxpool.Pool
	IR items.Repository
	JR jobs.Repository
	MR struct {
		quotations messages.Repository
		languages  messages.Repository
	}
	stripXML *bluemonday.Policy
}

// Init connects to all the app's resources and sets the config
func (app *App) Init(ctx context.Context) error {
	log.Info("Starting qftext")

	// Set the logging level
	ll, exists := os.LookupEnv("CCHC_LOGLEVEL")
	if !exists {
		ll = "info"
	}

	switch ll {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	log.Info("Connecting to the database")
	dbstr, exists := os.LookupEnv("CCHC_DBSTR")
	if !exists {
		return errors.New("Database connection string not set: use CCHC_DBSTR environment variable")
	}
	db, err := db.Connect(ctx, dbstr)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %w", err)
	}
	app.DB = db
	log.Info("Successfully connected to the database")

	app.IR = items.NewItemRepo(app.DB)
	app.JR = jobs.NewJobsRepo(app.DB)

	// Function to strip HTML/XML
	app.stripXML = bluemonday.StrictPolicy()

	// Connect to RabbitMQ and set up the queue
	mqstr, exists := os.LookupEnv("CCHC_MQSTR")
	if !exists {
		return errors.New("Message broker connection string not set: use CCHC_MQSTR environment variable")
	}
	log.Info("Attempting to connect to the message broker")

	quotations, err := messages.NewRabbitMQ(ctx, mqstr, "fulltext-quotations", 100)
	if err != nil {
		return fmt.Errorf("Error connecting to message broker: %w", err)
	}
	app.MR.quotations = quotations
	log.Info("Successfully connected to the message broker")

	languages, err := messages.NewRabbitMQ(ctx, mqstr, "fulltext-languages", 100)
	if err != nil {
		return fmt.Errorf("Error connecting to message broker: %w", err)
	}
	app.MR.languages = languages
	log.Info("Successfully connected to the message broker")

	return nil
}

// Shutdown closes the app's resources
func (app *App) Shutdown() {
	app.DB.Close()
	log.Info("Closed the connection to the database")
	app.MR.quotations.Close()
	app.MR.languages.Close()
	log.Info("Closed the connection to the message broker")
	log.Info("Stopped qftext")
}
