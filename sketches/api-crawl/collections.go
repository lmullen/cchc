package main

import (
	"fmt"
	"net/url"
	"strings"
)

// Starting place for the API
const apiBase = "https://www.loc.gov"

// Stuff we don't want or need from the API, which reduces the response size
var removeFromResponse = []string{
	"aka", "breadcrumbs", "categories", "content", "content_is_post",
	"expert_resources", "facet_trail", "facet_views", "facets", "featured_items",
	"form_facets", "legacy-url", "next", "next_sibling", "options",
	"original_formats", "pages", "partof", "previous", "previous_sibling",
	"research-centers", "shards", "site_type", "subjects", "timeline_1852_1880",
	"timeline_1881_1900", "timeline_1901_1925", "timestamp", "topics", "views",
}

// apiOptions set various parameters for the requests to the API
var apiOptions = url.Values{
	"at!": []string{strings.Join(removeFromResponse, ",")},
	"c":   []string{"250"},
	"fa":  []string{"online-format:online text"},
	"fo":  []string{"json"},
	"st":  []string{"list"},
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
