package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

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
	Consumer amqp.Delivery
}

// The App type shares access to the database and other resources.
type App struct {
	DB       *sql.DB
	Config   *Config
	Client   *http.Client
	MessageQ *amqp.Connection
	Limiters struct {
		Newspapers  ratelimit.Limiter
		Items       ratelimit.Limiter
		Collections ratelimit.Limiter
	}
	ItemMetadataQ Queue
}

// Init creates a new App and connects to the database or returns an error
func (app *App) Init() error {

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

	// Connect to the database and initialize it.
	log.Infof("Connecting to the %v database", app.Config.dbname)
	dbconstr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		app.Config.dbhost, app.Config.dbport, app.Config.dbname, app.Config.dbuser,
		app.Config.dbpass, app.Config.dbssl)
	db, err := sql.Open("pgx", dbconstr)
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("Failed to ping database: %w", err)
	}

	app.DB = db
	err = app.DBCreateSchema()
	if err != nil {
		return fmt.Errorf("Failed to create database schema: %w", err)
	}

	// Connect to RabbitMQ and set up the queues
	log.Info("Connecting to the message queue")
	qconnstr := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		app.Config.quser, app.Config.qpass, app.Config.qhost, app.Config.qport)
	rabbit, err := amqp.Dial(qconnstr)
	if err != nil {
		return fmt.Errorf("Failed to connect to message queue: %w", err)
	}
	app.MessageQ = rabbit
	ch, err := rabbit.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel on message queue: %w", err)
	}
	app.ItemMetadataQ.Channel = ch
	q, err := ch.QueueDeclare("items-metadata", true, false, false, false,
		amqp.Table{"x-max-length": 10000000, "x-queue-mode": "lazy"})
	if err != nil {
		return fmt.Errorf("Failed to declare a queue: %w", err)
	}
	app.ItemMetadataQ.Queue = &q
	consumer, err := ch.Consume(q.Name, "item-metadata-consumer",
		false, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to register a consumer", err)
	}

	// Set up a client to use for all HTTP requests
	app.Client = &http.Client{
		Timeout: apiTimeout * time.Second,
	}

	// Create rate limiters for different endpoints. Rate limits documentation:
	// https://www.loc.gov/apis/json-and-yaml/
	// TODO: The subtractions here represent a buffer from the officially presented
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
	err := app.DB.Close()
	if err != nil {
		log.Error("Failed to close the connection to the database:", err)
	}
	log.Info("Closing the connection to the message queue")
	err = app.MessageQ.Close()
	if err != nil {
		log.Error("Failed to close the connection to the message queue: ", err)
	}
}
