package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
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

// Queue holds all the items related to sending/receiving on a message queue
type Queue struct {
	Channel  *amqp.Channel
	Queue    *amqp.Queue
	Consumer <-chan amqp.Delivery
}

// The App type shares access to the database and other resources.
type App struct {
	DB            *pgxpool.Pool
	Config        *Config
	MessageBroker *amqp.Connection
	DocumentsQ    Queue
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

	db, err := pgxpool.Connect(ctx, app.Config.dbstr)
	if err != nil {
		return fmt.Errorf("Failed to dial the database: %w", err)
	}
	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("Failed to ping the database: %w", err)
	}
	app.DB = db
	log.Info("Connected to the database successfully")

	// Retry policy
	policy := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10)

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
	err = ch.Qos(100, 0, true)
	if err != nil {
		log.Fatal("Failed to set prefetch on the message broker: ", err)
	}
	app.DocumentsQ.Channel = ch
	dle, dlq, dlk := "failed-items-metadata", "dead-letter-queue", "dead-letter-key"
	err = ch.ExchangeDeclare(dle, "fanout", true, false, false, false, amqp.Table{})
	if err != nil {
		return fmt.Errorf("Failed to declare the dead letter exchange: %w", err)
	}
	_, err = ch.QueueDeclare(dlq, true, false, false, false, amqp.Table{})
	if err != nil {
		return fmt.Errorf("Failed to declare the dead letter exchange: %w", err)
	}
	err = ch.QueueBind(dlq, dlk, dle, false, amqp.Table{})
	if err != nil {
		return fmt.Errorf("Failed to bind dead letter queue and exchange: %w", err)
	}
	q, err := ch.QueueDeclare("documents", true, false, false, false,
		amqp.Table{
			"x-max-length":           10000000,
			"x-queue-mode":           "lazy",
			"x-dead-letter-exchange": dle,
		})
	if err != nil {
		return fmt.Errorf("Failed to declare a queue: %w", err)
	}
	app.DocumentsQ.Queue = &q
	consumer, err := ch.Consume(q.Name, "documents-consumer",
		false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Failed to register a channel consumer: %w", err)
	}
	app.DocumentsQ.Consumer = consumer
	log.Info("Connected to the message broker successfully")

	return nil
}

// Shutdown closes the connection to the database.
func (app *App) Shutdown() {
	app.DB.Close()
	log.Info("Closed the connection to the database")
	err := app.MessageBroker.Close()
	if err != nil {
		log.Error("Failed to close the connection to the message queue: ", err)
	} else {
		log.Info("Closed the connection to the message broker successfully")
	}
	log.Info("Shutdown the predictor")
}
