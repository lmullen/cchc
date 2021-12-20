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
	apiTimeout = 15 // The timeout limit for API requests in seconds
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
func (app *App) Init() error {
	log.Info("Starting the item metadata fetcher")

	app.Config = &Config{}

	ll, ok := os.LookupEnv("CCHC_LOGLEVEL")
	if !ok {
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

	app.Failures = make(map[string]time.Time)

	// Read the configuration from environment variables.
	dbstr, ok := os.LookupEnv("CCHC_DBSTR")
	if !ok {
		return errors.New("CCHC_DBSTR environment variable is not set")
	}
	app.Config.dbstr = dbstr

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := db.Connect(ctx, app.Config.dbstr)
	if err != nil {
		return err
	}
	app.DB = db
	app.ItemsRepo = items.NewItemRepo(db)
	log.Info("Connected to the database successfully")

	// Set up a client to use for all HTTP requests. It will automatically retry.
	rc := retryablehttp.NewClient()
	rc.RetryWaitMin = 3 * time.Second
	rc.RetryWaitMax = 20 * time.Second
	rc.RetryMax = 2
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
	app.DB.Close()
	log.Info("Closed the connection to the database successfully")
	log.Info("Shutdown the item metadata fetcher")
}
