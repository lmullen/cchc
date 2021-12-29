package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/lmullen/cchc/common/db"
	"github.com/lmullen/cchc/common/items"
	"github.com/lmullen/cchc/common/jobs"
	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4/pgxpool"
)

// The Config type stores configuration which is read from environment variables.
type Config struct {
	dbstr    string
	loglevel string
}

// The App type shares access to the database and other resources.
type App struct {
	DB        *pgxpool.Pool
	Config    *Config
	ItemsRepo items.Repository
	JobsRepo  jobs.Repository
}

// Init creates a new app and connects to the database or returns an error
func (app *App) Init(ctx context.Context) error {
	log.Info("Starting the language detector")

	// Set a timeout for getting the application set up
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	app.Config = &Config{}

	// Set the logging level
	ll, ok := os.LookupEnv("CCHC_LOGLEVEL")
	if !ok {
		ll = "info"
	}
	app.Config.loglevel = ll
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
	dbstr, ok := os.LookupEnv("CCHC_DBSTR")
	if !ok {
		return errors.New("CCHC_DBSTR environment variable is not set")
	}
	app.Config.dbstr = dbstr

	db, err := db.Connect(ctx, app.Config.dbstr, "cchc-language-detector")
	if err != nil {
		return err
	}
	app.DB = db
	app.ItemsRepo = items.NewItemRepo(db)
	app.JobsRepo = jobs.NewJobsRepo(db)
	log.Info("Connected to the database successfully")

	return nil
}

// Shutdown closes the connection to the database.
func (app *App) Shutdown() {
	log.Info("Closing the connection to the database")
	app.DB.Close()
	log.Info("Shutdown the language detector")
}
