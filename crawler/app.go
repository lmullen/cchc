package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/ratelimit"
)

// The Config type stores configuration which is read from environment variables.
type Config struct {
	dbhost string
	dbport string
	dbname string
	dbuser string
	dbpass string
	dbssl  string // SSL mode for the database connection
}

// The App type shares access to the database and other resources.
type App struct {
	DB                 *sql.DB
	Config             *Config
	Client             *http.Client
	NewspaperLimiter   ratelimit.Limiter
	ItemsLimiter       ratelimit.Limiter
	CollectionsLimiter ratelimit.Limiter
	CollectionsWG      *sync.WaitGroup
}

// getEnv either returns the value of an environment variable or, if that
// environment variables does not exist, returns the fallback value provided.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
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

	// Connect to the database then store the database in the struct.
	log.Infof("Connecting to the %v database", app.Config.dbname)
	constr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		app.Config.dbhost, app.Config.dbport, app.Config.dbname, app.Config.dbuser,
		app.Config.dbpass, app.Config.dbssl)
	db, err := sql.Open("pgx", constr)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	app.DB = db
	app.DBInit()

	app.Client = &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create rate limiters for different endpoints. Rate limits documentation:
	// https://www.loc.gov/apis/json-and-yaml/
	// TODO: The subtractions here represent a buffer from the officially presented
	// rate limits.
	il := ratelimit.New(200-20, ratelimit.Per(60*time.Second)) // 200 requests/minute
	app.ItemsLimiter = il

	cl := ratelimit.New(80-20, ratelimit.Per(60*time.Second)) // 80 requests/minute
	app.CollectionsLimiter = cl

	nl := ratelimit.New(20-4, ratelimit.Per(10*time.Second)) // 120 requests/minute
	app.NewspaperLimiter = nl

	app.CollectionsWG = &sync.WaitGroup{}

	return nil
}

// Shutdown closes the connection to the database.
func (app *App) Shutdown() {
	log.Info("Closing the connection to the database")
	err := app.DB.Close()
	if err != nil {
		log.Error(err)
	}
}

// Exit the entire program if we get an HTTP 429 error
// TODO: Would be better to wait and try again, but this works for now
func quitIfBlocked(code int) {
	if code == http.StatusTooManyRequests {
		app.Shutdown()
		log.Fatal("Quiting because rate limit exceeded")
	}
}
