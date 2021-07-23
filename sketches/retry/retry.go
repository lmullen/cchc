package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Set up a client to use for all HTTP requests. It will automatically retry.
	rc := retryablehttp.NewClient()
	rc.RetryWaitMin = 2 * time.Second
	rc.RetryWaitMax = 16 * time.Second
	rc.RetryMax = 3
	rc.HTTPClient.Timeout = 60 * time.Second
	// rc.Logger = nil
	rc.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
		log.WithFields(logrus.Fields{
			"attempt": attempt,
			"url":     req.URL,
		}).Debug("Fetching URL")
	}
	client := rc.StandardClient()

	var wg sync.WaitGroup
	wg.Add(2)

	// Test 5xx errors
	go func() {
		defer wg.Done()
		res, err := client.Get("https://httpstat.us/500")
		if err != nil {
			log.Error("Error: ", err)
		}
		log.Info("Result: ", res)
	}()

	// Test 429 errors
	go func() {
		defer wg.Done()
		res, err := client.Get("https://httpstat.us/429")
		if err != nil {
			log.Error("Error: ", err)
		}
		log.Info("Result: ", res)
	}()

	wg.Wait()

}
