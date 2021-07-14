package main

import (
	"net/http"
	"os"
	"strings"

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

// hasAPI checks whether the URL is a part of the LOC.gov API, as opposed to
// American Memory or some other digital collection.
func hasAPI(url string) bool {
	return strings.HasPrefix(url, apiBase)
}
