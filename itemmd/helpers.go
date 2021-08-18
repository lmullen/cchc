package main

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// getEnv either returns the value of an environment variable or, if that
// environment variables does not exist, returns the fallback value provided.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// Exit the entire program if we get an HTTP 429 error for making too many requests.
// TODO: Would be better to wait and try again, but this works for now
func quitIfBlocked(code int) {
	if code == http.StatusTooManyRequests {
		app.Shutdown()
		log.Fatal("Quiting because rate limit exceeded")
	}
}

// year takes an ISO8601 date string and figures out the year
func year(date string) sql.NullInt32 {
	var year sql.NullInt32
	if len(date) >= 4 {
		s := date[0:4]
		y, err := strconv.Atoi(s)
		if err == nil {
			year.Scan(y)
		}
	}
	return year
}
