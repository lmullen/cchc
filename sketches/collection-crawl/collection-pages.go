package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Starting place for the API
const apiBase = "https://www.loc.gov"

// Stuff we don't want or need from the API, which reduces the response size
var removeFromResponse = []string{
	"aka", "breadcrumbs", "browse", "categories", "content", "content_is_post",
	"expert_resources", "facet_trail", "facet_views", "facets", "featured_items",
	"form_facets", "legacy-url", "next", "next_sibling", "options",
	"original_formats", "pages", "partof", "previous", "previous_sibling",
	"research-centers", "shards", "site_type", "subjects", "timeline_1852_1880",
	"timeline_1881_1900", "timeline_1901_1925", "timestamp", "topics", "views",
}

// apiOptions set various parameters for the requests to the API
var apiCollectionOptions = url.Values{
	"at!": []string{strings.Join(removeFromResponse, ",")},
	"c":   []string{fmt.Sprint(itemsPerPage)},
	"fa":  []string{"online-format:online text"},
	"fo":  []string{"json"},
	"st":  []string{"list"},
}

// CollectionURL takes a slug for a particular collecction, plus the page in that
// collections results to fetch, and returns the URL for that page of that collection.
func CollectionURL(slug string, page int) string {
	urlBase := apiBase + "/collections/" + slug + "/"
	u, _ := url.Parse(urlBase)

	// Set the query to be the API options, then add the correct page of results
	q := apiCollectionOptions
	q.Set("sp", fmt.Sprint(page))
	u.RawQuery = q.Encode()

	return u.String()
}

// CollectionAPIPage is an object returned by querying a specific page of the
// collections endpoint of the LOC.gov API.
type CollectionAPIPage struct {
	// ContentIsPost bool `json:"content_is_post"`
	// Digitized int `json:"digitized"`
	// FormFacets    struct {
	// } `json:"form_facets"`
	Pagination struct {
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
	// Search  struct {
	// 	Dates       interface{} `json:"dates"`
	// 	FacetLimits string      `json:"facet_limits"`
	// 	Field       interface{} `json:"field"`
	// 	Hits        int         `json:"hits"`
	// 	In          string      `json:"in"`
	// 	Query       string      `json:"query"`
	// 	Recommended int         `json:"recommended"`
	// 	Site        struct {
	// 	} `json:"site"`
	// 	SortBy      string `json:"sort_by"`
	// 	Type        string `json:"type"`
	// 	UnionFacets string `json:"union_facets"`
	// 	URL         string `json:"url"`
	// } `json:"search"`
	Title string `json:"title"`
	Total int    `json:"total"`
}

// String prints the collection result
func (collection CollectionAPIPage) String() string {
	out := fmt.Sprintf("%s, page %v", collection.Title, collection.Pagination.Current)
	return out
}

func fetchCollectionResult(url string, client *http.Client, results chan<- CollectionAPIPage) {

	response, err := client.Get(url)
	if err != nil {
		log.Println(err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var result CollectionAPIPage

	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Println(err)
	}

	results <- result

	// If there is another page of results, go fetch it in a goroutine. Otherwise,
	// close the channel so we know we are done with results.
	if result.Pagination.Next != "" {
		url := CollectionURL(collectionSlug, result.Pagination.Current+1)
		go fetchCollectionResult(url, client, results)
	} else {
		close(results)
	}

}
