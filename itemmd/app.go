package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/lmullen/cchc/common/db"
	"github.com/lmullen/cchc/common/items"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/ratelimit"
)

// Configuration options that aren't worth exposing as environment variables
const (
	apiTimeout = 10 // The timeout limit for API requests in seconds
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
	Client    *http.Client
	ItemsRepo items.Repository
	Limiters  struct {
		Newspapers  ratelimit.Limiter
		Items       ratelimit.Limiter
		Collections ratelimit.Limiter
	}
	Failures map[string]time.Time
}

// Init creates a new app and connects to the database or returns an error
func (app *App) Init(ctx context.Context) error {
	log.Info("Starting the item metadata fetcher")

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

	// Record items that we failed to fetch and when so we don't get stuck fetching them
	app.Failures = make(map[string]time.Time)

	// Connect to the database and create the various repositories needed
	dbstr, ok := os.LookupEnv("CCHC_DBSTR")
	if !ok {
		return errors.New("CCHC_DBSTR environment variable is not set")
	}
	app.Config.dbstr = dbstr

	db, err := db.Connect(ctx, app.Config.dbstr)
	if err != nil {
		return err
	}
	app.DB = db
	app.ItemsRepo = items.NewItemRepo(db)
	log.Info("Connected to the database successfully")

	// Set up a client to use for all HTTP requests. It will automatically retry.
	rc := retryablehttp.NewClient()
	rc.RetryWaitMin = 2 * time.Second
	rc.RetryWaitMax = 5 * time.Second
	rc.RetryMax = 3
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
	log.Info("Closing the connection to the database")
	app.DB.Close()
	log.Info("Shutdown the item metadata fetcher")
}
