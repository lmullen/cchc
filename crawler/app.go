package main

import (
	"database/sql"
	"fmt"
	"net/http"
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
	dbhost   string
	dbport   string
	dbname   string
	dbuser   string
	dbpass   string
	dbssl    string // SSL mode for the database connection
	qhost    string
	quser    string
	qpass    string
	qport    string
	loglevel string
}

// Queue holds all the items related to sending/receiving on a message queue
type Queue struct {
	Channel  *amqp.Channel
	Queue    *amqp.Queue
	Consumer <-chan amqp.Delivery
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

	// Read the configuration from environment variables. The `getEnv()` function
	// will provide a default.
	app.Config.dbhost = getEnv("CCHC_DBHOST", "localhost")
	app.Config.dbport = getEnv("CCHC_DBPORT", "5432")
	app.Config.dbname = getEnv("CCHC_DBNAME", "cchc")
	app.Config.dbuser = getEnv("CCHC_DBUSER", "lmullen")
	app.Config.dbpass = getEnv("CCHC_DBPASS", "")
	app.Config.qhost = getEnv("CCHC_QHOST", "localhost")
	app.Config.qport = getEnv("CCHC_QPORT", "5672")
	app.Config.quser = getEnv("CCHC_QUSER", "cchc")
	app.Config.qpass = getEnv("CCHC_QPASS", "")
	app.Config.loglevel = getEnv("CCHC_LOGLEVEL", "warn")

	// Set the logging level
	switch app.Config.loglevel {
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
	dbconstr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		app.Config.dbhost, app.Config.dbport, app.Config.dbname, app.Config.dbuser,
		app.Config.dbpass, app.Config.dbssl)
	var db *sql.DB
	dbConnect := func() error {
		d, err := sql.Open("pgx", dbconstr)
		if err != nil {
			return fmt.Errorf("Failed to dial the database: %w", err)
		}
		if err := d.Ping(); err != nil {
			return fmt.Errorf("Failed to ping the database: %w", err)
		}
		db = d
		return nil
	}
	log.Info("Attempting to connect to the database")
	err := backoff.Retry(dbConnect, policy)
	if err != nil {
		return fmt.Errorf("Failed to connect to the database: %w", err)
	}

	app.DB = db
	err = app.DBCreateSchema()
	if err != nil {
		return fmt.Errorf("Failed to create database schema: %w", err)
	}
	log.Info("Connected to the database successfully")

	// Connect to RabbitMQ and set up the queues. Try to connect multiple times
	qconnstr := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		app.Config.quser, app.Config.qpass, app.Config.qhost, app.Config.qport)
	var rabbit *amqp.Connection
	mqConnect := func() error {
		r, err := amqp.Dial(qconnstr)
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
	err = ch.Qos(40, 0, true)
	if err != nil {
		log.Fatal("Failed to set prefetch on the message broker: ", err)
	}
	app.ItemMetadataQ.Channel = ch
	q, err := ch.QueueDeclare("items-metadata", true, false, false, false,
		amqp.Table{
			"x-max-length":           10000000,
			"x-queue-mode":           "lazy",
			"x-dead-letter-exchange": "failed-items-metadata",
		})
	if err != nil {
		return fmt.Errorf("Failed to declare a queue: %w", err)
	}
	app.ItemMetadataQ.Queue = &q
	consumer, err := ch.Consume(q.Name, "item-metadata-consumer",
		false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Failed to register a channel consumer: %w", err)
	}
	app.ItemMetadataQ.Consumer = consumer
	log.Info("Connected to the message broker successfully")

	// Set up a client to use for all HTTP requests. It will automatically retry.
	rc := retryablehttp.NewClient()
	rc.RetryWaitMin = 10 * time.Second
	rc.RetryWaitMax = 10 * time.Minute
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
