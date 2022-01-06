package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/go-retryablehttp"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/ratelimit"
)

// The Config type stores configuration which is read from environment variables.
type Config struct {
	dbstr    string
	loglevel string
}

// The App type shares access to the database and other resources.
type App struct {
	DB       *sql.DB
	Config   *Config
	Client   *http.Client
	Limiters struct {
		Newspapers  ratelimit.Limiter
		Items       ratelimit.Limiter
		Collections ratelimit.Limiter
	}
}

// Init creates a new App and connects to the database or returns an error
func (app *App) Init() error {
	log.Info("Starting the LOC.gov API crawler")

	app.Config = &Config{}

	// Read the configuration from environment variables.
	dbstr, ok := os.LookupEnv("CCHC_DBSTR")
	if !ok {
		return errors.New("CCHC_DBSTR environment variable is not set")
	}
	app.Config.dbstr = dbstr

	ll, ok := os.LookupEnv("CCHC_LOGLEVEL")
	if !ok {
		return errors.New("CCHC_LOGLEVEL environment variable is not set")
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

	// Set a policy for backoffs
	policy := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10)

	// Connect to the database and initialize it.
	var db *sql.DB
	dbConnect := func() error {
		d, err := sql.Open("pgx", app.Config.dbstr)
		if err != nil {
			return fmt.Errorf("Failed to dial the database: %w", err)
		}
		if err := d.Ping(); err != nil {
			return fmt.Errorf("Failed to ping the database: %w", err)
		}
		db = d
		return nil
	}
	log.Infof("Attempting to connect to the database")
	err := backoff.Retry(dbConnect, policy)
	if err != nil {
		return fmt.Errorf("Failed to connect to the database: %w", err)
	}

	app.DB = db
	log.Info("Connected to the database successfully")

	// Set up a client to use for all HTTP requests. It will automatically retry.
	rc := retryablehttp.NewClient()
	rc.RetryWaitMin = 10 * time.Second
	rc.RetryWaitMax = 2 * time.Minute
	rc.RetryMax = 6
	rc.HTTPClient.Timeout = apiTimeout * time.Second
	rc.Logger = nil
	// This will log all HTTP requests made, which is not desirable.
	// rc.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
	// 	log.WithFields(logrus.Fields{
	// 		"attempt": attempt,
	// 		"url":    req.URL,
	// 	}).Debug("Fetching URL")
	// }
	app.Client = rc.StandardClient()

	// Create rate limiters for different endpoints. Rate limits documentation:
	// https://www.loc.gov/apis/json-and-yaml/
	// The subtractions here represent a buffer from the officially presented
	// rate limits.
	il := ratelimit.New(200-20, ratelimit.Per(60*time.Second)) // 200 requests/minute
	app.Limiters.Items = il

	cl := ratelimit.New(80-20, ratelimit.Per(60*time.Second)) // 80 requests/minute
	app.Limiters.Collections = cl

	nl := ratelimit.New(20-4, ratelimit.Per(10*time.Second)) // 120 requests/minute
	app.Limiters.Newspapers = nl

	return nil
}

// Shutdown closes the connection to the database.
func (app *App) Shutdown() {
	err := app.DB.Close()
	if err != nil {
		log.Error("Failed to close the connection to the database:", err)
	} else {
		log.Info("Closed the connection to the database successfully")
	}
	log.Info("Shutdown the LOC.gov API crawler")
}
