package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
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

// The App type shares access to the database.
type App struct {
	DB     *sql.DB
	Config *Config
	Client *http.Client
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
	log.Printf("Connecting to the %v database.\n", app.Config.dbname)
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

	app.Client = &http.Client{}

	return nil
}

// Shutdown closes the connection to the database and shutsdown the server.
func (app *App) Shutdown() error {
	log.Println("Closing the connection to the database.")
	err := app.DB.Close()
	if err != nil {
		return err
	}
	return nil
}
