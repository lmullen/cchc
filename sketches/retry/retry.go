package main

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Set up a client to use for all HTTP requests. It will automatically retry.
	rc := retryablehttp.NewClient()
	rc.RetryWaitMin = 2 * time.Second
	rc.RetryWaitMax = 2 * time.Minute
	rc.RetryMax = 6
	rc.HTTPClient.Timeout = 60 * time.Second
	// rc.Logger = nil
	rc.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
		log.WithFields(logrus.Fields{
			"attempt": attempt,
			"url":     req.URL,
		}).Debug("Fetching URL")
	}
	client := rc.StandardClient()

	res, err := client.Get("https://httpstat.us/500")
	if err != nil {
		log.Error("Error: ", err)
	}
	log.Info("Result: ", res)
}
