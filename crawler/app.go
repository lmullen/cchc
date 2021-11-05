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
	"github.com/streadway/amqp"

	"github.com/hashicorp/go-retryablehttp"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/ratelimit"
)

// The Config type stores configuration which is read from environment variables.
type Config struct {
	dbstr    string
	mqstr    string
	loglevel string
}

// Queue holds all the items related to sending on a message queue
type Queue struct {
	Channel *amqp.Channel
	Queue   *amqp.Queue
}

// The App type shares access to the database and other resources.
type App struct {
	DB            *sql.DB
	Config        *Config
	Client        *http.Client
	MessageBroker *amqp.Connection
	Limiters      struct {
		Newspapers  ratelimit.Limiter
		Items       ratelimit.Limiter
		Collections ratelimit.Limiter
	}
	ItemMetadataQ Queue
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

	mqstr, ok := os.LookupEnv("CCHC_MQSTR")
	if !ok {
		return errors.New("CCHC_MQSTR environment variable is not set")
	}
	app.Config.mqstr = mqstr

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

	// Connect to RabbitMQ and set up the queues. Try to connect multiple times
	var rabbit *amqp.Connection
	mqConnect := func() error {
		r, err := amqp.Dial(app.Config.mqstr)
		if err != nil {
			return err
		}
		rabbit = r
		return nil
	}

	log.Info("Attempting to connect to the message broker")
	err = backoff.Retry(mqConnect, policy)
	if err != nil {
		return fmt.Errorf("Failed to connect to message broker: %w", err)
	}
	app.MessageBroker = rabbit
	ch, err := rabbit.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel on message broker: %w", err)
	}
	err = ch.Qos(64, 0, true)
	if err != nil {
		log.Fatal("Failed to set prefetch on the message broker: ", err)
	}
	app.ItemMetadataQ.Channel = ch
	q, err := ch.QueueDeclare("items-metadata", true, false, false, false,
		amqp.Table{
			"x-queue-mode":           "lazy",
			"x-dead-letter-exchange": "dead-letter-exchange",
		})
	if err != nil {
		return fmt.Errorf("Failed to declare a queue: %w", err)
	}
	app.ItemMetadataQ.Queue = &q
	log.Info("Connected to the message broker successfully")

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
	err = app.MessageBroker.Close()
	if err != nil {
		log.Error("Failed to close the connection to the message queue: ", err)
	} else {
		log.Info("Closed the connection to the message broker successfully")
	}
	log.Info("Shutdown the LOC.gov API crawler")
}
