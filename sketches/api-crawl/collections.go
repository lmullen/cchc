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

// CollectionResult is an object returned by querying the collections endpoint of
// of the LOC.gov API.
type CollectionResult struct {
	ContentIsPost bool `json:"content_is_post"`
	Digitized     int  `json:"digitized"`
	FormFacets    struct {
	} `json:"form_facets"`
	Pagination struct {
		Current  int    `json:"current"`
		First    string `json:"first"`
		From     int    `json:"from"`
		Last     string `json:"last"`
		Next     string `json:"next"`
		Of       int    `json:"of"`
		PageList []struct {
			Number int    `json:"number"`
			URL    string `json:"url"`
		} `json:"page_list"`
		Perpage        int    `json:"perpage"`
		PerpageOptions []int  `json:"perpage_options"`
		Previous       string `json:"previous"`
		Results        string `json:"results"`
		To             int    `json:"to"`
		Total          int    `json:"total"`
	} `json:"pagination"`
	Results []ItemResult `json:"results"`
	Search  struct {
		Dates       interface{} `json:"dates"`
		FacetLimits string      `json:"facet_limits"`
		Field       interface{} `json:"field"`
		Hits        int         `json:"hits"`
		In          string      `json:"in"`
		Query       string      `json:"query"`
		Recommended int         `json:"recommended"`
		Site        struct {
		} `json:"site"`
		SortBy      string `json:"sort_by"`
		Type        string `json:"type"`
		UnionFacets string `json:"union_facets"`
		URL         string `json:"url"`
	} `json:"search"`
	Title string `json:"title"`
	Total int    `json:"total"`
}
