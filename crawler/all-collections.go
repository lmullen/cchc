package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// FetchAllCollections gets all the digital collections that match the query.
//
// TODO: Note that we are not worrying about pagination since there are currently
// fewer digital collections that are available and that we care about than the
// pagination limit.
func FetchAllCollections() ([]Collection, error) {

	// Rate limiter
	app.Limiters.Collections.Take()

	// Build the URL with the correct query
	u, _ := url.Parse(apiBase + "/collections/")

	apiAllCollectionOptions := url.Values{
		"at!": []string{strings.Join(removeFromResponse, ",")},
		"c":   []string{fmt.Sprint(apiItemsPerPage)},
		"fa":  []string{"subject_topic:american history"}, // TODO: Consider removing subject limit
		"fo":  []string{"json"},
	}
	u.RawQuery = apiAllCollectionOptions.Encode()
	url := u.String()

	log.Info("Fetching all digital collections")
	log.WithField("url", url).Debug("URL for digital collections")
	response, err := app.Client.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"http_error": response.Status,
			"http_code":  response.StatusCode,
			"url":        url,
		}).Warn("HTTP error when fetching from API")
		quitIfBlocked(response.StatusCode)
		return nil, fmt.Errorf("HTTP error: %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading HTTP response body: %w", err)
	}

	var result CollectionsList

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling collections list: %w", err)
	}

	return result.Results, nil

}

// CollectionsList is the API response with all collections, containing an array
// of Collections. The API returns many more fields than this, but they are
// ignored when unmarshalling the JSON.
type CollectionsList struct {
	// Pagination struct {
	// 	Current  int         `json:"current"`
	// 	First    interface{} `json:"first"`
	// 	From     int         `json:"from"`
	// 	Last     string      `json:"last"`
	// 	Next     string      `json:"next"`
	// 	Of       int         `json:"of"`
	// 	PageList []struct {
	// 		Number int         `json:"number"`
	// 		URL    interface{} `json:"url"`
	// 	} `json:"page_list"`
	// 	Perpage        int         `json:"perpage"`
	// 	PerpageOptions []int       `json:"perpage_options"`
	// 	Previous       interface{} `json:"previous"`
	// 	Results        string      `json:"results"`
	// 	To             int         `json:"to"`
	// 	Total          int         `json:"total"`
	// } `json:"pagination"`
	Results []Collection `json:"results"`
}
