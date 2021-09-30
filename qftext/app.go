package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lmullen/cchc/common/db"
	"github.com/lmullen/cchc/common/items"
	"github.com/lmullen/cchc/common/jobs"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
)

// App holds resources and config
type App struct {
	DB       *pgxpool.Pool
	IR       items.Repository
	JR       jobs.Repository
	stripXML *bluemonday.Policy
}

// Init connects to all the app's resources and sets the config
func (app *App) Init() error {
	log.Info("Creating jobs for the prediction model from items with full text")

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

	log.Info("Connecting to the database")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	db, err := db.Connect(ctx)
	if err != nil {
		return err
	}
	app.DB = db
	log.Info("Successfully connected to the database")

	app.IR = items.NewItemRepo(app.DB)
	app.JR = jobs.NewJobsRepo(app.DB)

	app.stripXML = bluemonday.StrictPolicy()

	return nil
}

// Shutdown closes the app's resources
func (app *App) Shutdown() {
	app.DB.Close()
	log.Info("Closed the connection to the database")
	log.Info("Stopped process to create jobs for predictions from full text")
}
