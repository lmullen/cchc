package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/lmullen/cchc/common/db"
	"github.com/lmullen/cchc/common/items"
	"github.com/lmullen/cchc/common/jobs"
	"github.com/lmullen/cchc/common/results"
	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Configuration options that aren't worth exposing as environment variables
const (
	apiTimeout = 60 // The timeout limit for API requests in seconds
)

// The Config type stores configuration which is read from environment variables.
type Config struct {
	dbstr    string
	loglevel string
}

// The App type shares access to the database and other resources.
type App struct {
	DB          *pgxpool.Pool
	Config      *Config
	ItemsRepo   items.Repository
	ResultsRepo results.Repository
	JobsRepo    jobs.Repository
}

// Init creates a new app and connects to the database or returns an error
func (app *App) Init() error {
	log.Info("Starting the prediction modeler fetcher")

	app.Config = &Config{}

	// Set a timeout for getting the application set up
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ll, exists := os.LookupEnv("CCHC_LOGLEVEL")
	if !exists {
		ll = "info"
	}
	app.Config.loglevel = ll

	// Set the logging level
	switch app.Config.loglevel {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	}

	// Connect to the database and create the various repositories needed
	dbstr, exists := os.LookupEnv("CCHC_DBSTR")
	if !exists {
		return errors.New("CCHC_DBSTR environment variable is not set")
	}
	app.Config.dbstr = dbstr

	db, err := db.Connect(ctx, app.Config.dbstr, "cchc-predictor-aggregator")
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %w", err)
	}
	app.DB = db
	app.ItemsRepo = items.NewItemRepo(db)
	app.JobsRepo = jobs.NewJobsRepo(db)
	app.ResultsRepo = results.NewRepo(db)
	log.Info("Connected to the database successfully")

	// Initialize the results repo
	res := results.NewRepo(db)
	app.ResultsRepo = res

	// Initialize the jobs repo
	jobs := jobs.NewJobsRepo(db)
	app.JobsRepo = jobs

	return nil
}

// Shutdown closes the connection to the database.
func (app *App) Shutdown() {
	log.Info("Closing the connection to the database")
	app.DB.Close()
	log.Info("Shutdown the prediction model fetcher")
}
