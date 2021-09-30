package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lmullen/cchc/common/db"
	"github.com/lmullen/cchc/common/items"
	"github.com/lmullen/cchc/common/jobs"
	"github.com/lmullen/cchc/common/messages"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// App holds resources and config
type App struct {
	DB       *pgxpool.Pool
	IR       items.Repository
	JR       jobs.Repository
	MR       messages.Repository
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

	// Function to strip HTML/XML
	app.stripXML = bluemonday.StrictPolicy()

	// Connect to RabbitMQ and set up the queue

	mqstr, exists := os.LookupEnv("CCHC_MQSTR")
	if !exists {
		return errors.New("Message broker connection string not set; use CCHC_MQSTR environment variable")
	}
	log.Info("Attempting to connect to the message broker")
	rabbit, err := amqp.Dial(mqstr)
	if err != nil {
		return fmt.Errorf("Error connecting to message broker: %w", err)
	}

	ch, err := rabbit.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel on message broker: %w", err)
	}
	err = ch.Qos(8, 0, true)
	if err != nil {
		return fmt.Errorf("Failed to set prefetch on the message broker: %w", err)
	}
	dle, dlq, dlk := "failed-fulltext", "dead-letter-queue", "dead-letter-key"
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
	q, err := ch.QueueDeclare("fulltext-predict", true, false, false, false,
		amqp.Table{
			// "x-max-length":           10000000,
			"x-queue-mode":           "lazy",
			"x-dead-letter-exchange": dle,
		})
	if err != nil {
		return fmt.Errorf("Failed to declare a queue: %w", err)
	}

	app.MR = messages.NewMessageRepo(ch, &q, nil)
	log.Info("Successfully connected to the message broker")

	return nil
}

// Shutdown closes the app's resources
func (app *App) Shutdown() {
	app.DB.Close()
	log.Info("Closed the connection to the database")
	log.Info("Stopped process to create jobs for predictions from full text")
}
