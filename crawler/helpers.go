package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Exit the entire program if we get an HTTP 429 error for making too many requests.
// TODO: Would be better to wait and try again, but this works for now
func quitIfBlocked(code int) {
	if code == http.StatusTooManyRequests {
		app.Shutdown()
		log.Fatal("Quiting because rate limit exceeded")
	}
}

// hasAPI checks whether the URL is a part of the LOC.gov API, as opposed to
// American Memory or some other digital collection.
func hasAPI(url string) bool {
	return strings.HasPrefix(url, apiBase)
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
