package main

import (
	"fmt"
	"net/url"
)

// Starting place for the API
const apiBase = "https://www.loc.gov"

// apiOptions set various parameters for the requests to the API
var apiOptions = url.Values{
	"c":  []string{"250"},
	"fa": []string{"online-format:online text"},
	"fo": []string{"json"},
	"st": []string{"list"},
}

// CollectionURL takes a slug for a particular collecction, plus the page in that
// collections results to fetch, and returns the URL for that page of that collection.
func CollectionURL(slug string, page int) (string, error) {
	urlBase := apiBase + "/collections/" + slug + "/"
	u, err := url.Parse(urlBase)

	if err != nil {
		return "", err
	}

	// Set the query to be the API options, then add the correct page of results
	q := apiOptions
	q.Set("sp", fmt.Sprint(page))
	u.RawQuery = q.Encode()

	return u.String(), nil
}
