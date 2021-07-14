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

// CollectionPageURL takes the URL to the items for a particular collecction, plus
// the page in that collections results to fetch, and returns the URL for that
// page of that collection.
func collectionPageURL(itemsURL string, page int) string {
	u, _ := url.Parse(itemsURL)

	// Set the query to be the API options, then add the correct page of results
	q := url.Values{
		"at!": []string{strings.Join(removeFromResponse, ",")},
		"c":   []string{fmt.Sprint(apiItemsPerPage)},
		"fo":  []string{"json"},
		"st":  []string{"list"},
		// "fa":  []string{"online-format:online text"}, // Not sure if this is a good query
	}
	q.Set("sp", fmt.Sprint(page))
	u.RawQuery = q.Encode()

	return u.String()
}

// FetchCollectionItems gets the items associated with each collection
func (c Collection) FetchCollectionItems(page int, results chan<- CollectionAPIPage) {

	defer app.CollectionsWG.Done()

	url := collectionPageURL(c.ItemsURL, page)

	// Skip if it isn't a part of the LOC.gov API
	if !hasAPI(url) {
		return
	}

	// Limit the rate
	app.Limiters.Collections.Take()

	response, err := app.Client.Get(url)
	if err != nil {
		log.Warn(err)
		return
	}

	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"http_error": response.Status,
			"http_code":  response.StatusCode,
			"url":        url,
		}).Warn("HTTP error when fetching from API")
		quitIfBlocked(response.StatusCode)
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Warn("Error reading HTTP response body: ", err)
		return
	}

	var result CollectionAPIPage

	err = json.Unmarshal(data, &result)
	if err != nil {
		log.WithFields(log.Fields{
			"url":           url,
			"parsing_error": err,
		}).Warn("Error parsing JSON")
		return // Quit early in the hopes of not messing up other go routines
	}

	log.Info("Fetched ", result)

	// Save the collectionID for creating a relation in the database
	result.CollectionID = c.ID

	results <- result

	// If there is another page of results, go fetch it.
	if result.Pagination.Next != "" {
		app.CollectionsWG.Add(1)
		c.FetchCollectionItems(result.Pagination.Current+1, results)
	}

}

// CollectionAPIPage is an object returned by querying a specific page of the
// collections endpoint of the LOC.gov API. Other fields are returned by the
// API but are ignored when parsing.
type CollectionAPIPage struct {
	CollectionID string // This is stored but is not returned as part of the API
	Pagination   struct {
		Current int `json:"current"`
		// First   string `json:"first"`
		// From    int    `json:"from"`
		// Last    string `json:"last"`
		Next string `json:"next"`
		// Of      int    `json:"of"`
		// PageList []struct {
		// 	Number int    `json:"number"`
		// 	URL    string `json:"url"`
		// } `json:"page_list"`
		// Perpage        int    `json:"perpage"`
		// PerpageOptions []int  `json:"perpage_options"`
		// Previous       string `json:"previous"`
		// Results        string `json:"results"`
		// To             int    `json:"to"`
		// Total          int    `json:"total"`
	} `json:"pagination"`
	Results []ItemResult `json:"results"`
	Title   string       `json:"title"`
}

// String prints the collection result
func (c CollectionAPIPage) String() string {
	out := fmt.Sprintf("%s, page %v", c.Title, c.Pagination.Current)
	return out
}
