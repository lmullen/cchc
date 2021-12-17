package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/lmullen/cchc/common/db"
	"github.com/lmullen/cchc/common/jobs"
	"github.com/lmullen/cchc/common/messages"
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
	mqstr    string
	loglevel string
}

// The App type shares access to the database and other resources.
type App struct {
	DB          *pgxpool.Pool
	Config      *Config
	ResultsRepo results.Repository
	JobsRepo    jobs.Repository
	MsgRepo     messages.Repository
}

// Init creates a new app and connects to the database or returns an error
func (app *App) Init() error {
	log.Info("Starting the prediction modeler fetcher")

	app.Config = &Config{}

	// Read the configuration from environment variables
	dbstr, exists := os.LookupEnv("CCHC_DBSTR")
	if !exists {
		return errors.New("Database connection string not set as an environment variable")
	}
	app.Config.dbstr = dbstr

	mqstr, exists := os.LookupEnv("CCHC_MQSTR")
	if !exists {
		return errors.New("Message queue connection string not set as an environment variable")
	}
	app.Config.mqstr = mqstr

	ll, exists := os.LookupEnv("CCHC_LOGLEVEL")
	if !exists {
		ll = "warn"
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
	default:
		log.SetLevel(log.WarnLevel)
	}

	// Connect to the database and initialize it.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := db.Connect(ctx, app.Config.dbstr)
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %w", err)
	}
	app.DB = db
	log.Info("Connected to the database successfully")

	log.Info("Attempting to connect to the message broker")
	rabbit, err := messages.NewRabbitMQ(ctx, app.Config.mqstr, "fulltext-predict", 100)
	if err != nil {
		return fmt.Errorf("Error connecting to message broker: %w", err)
	}
	app.MsgRepo = rabbit
	log.Info("Connected to the message broker successfully")

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
	app.DB.Close()
	log.Info("Closed the connection to the database")
	err := app.MsgRepo.Close()
	if err != nil {
		log.Error("Failed to close the connection to the message queue: ", err)
	} else {
		log.Info("Closed the connection to the message broker successfully")
	}
	log.Info("Shutdown the predictor")
}
