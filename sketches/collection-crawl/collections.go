package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var apiAllCollectionOptions = url.Values{
	"at!": []string{strings.Join(removeFromResponse, ",")},
	"c":   []string{"250"},
	"fa":  []string{"subject_topic:american history"},
	"fo":  []string{"json"},
}

// FetchAllCollections gets all the digital collections that match the query
// TODO: Note that we are not worrying about pagination since there are currently
// fewer digital collections that are available and that we care about than the
// pagination limit.
func FetchAllCollections(client *http.Client) ([]CollectionMetadata, error) {
	urlBase := apiBase + "/collections/"
	u, _ := url.Parse(urlBase)

	// Set the query to have the API options
	q := apiAllCollectionOptions
	u.RawQuery = q.Encode()

	url := u.String()

	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result CollectionsList

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result.Results, nil

	// https: //www.loc.gov/collections/?fa=subject_topic:american+history&fo=json&c=250
}

// CollectionsList contains an array of descriptions of Collection Metadata
type CollectionsList struct {
	// Browse struct {
	// 	Facets struct {
	// 		Include []struct {
	// 			Field string `json:"field"`
	// 			Label string `json:"label"`
	// 		} `json:"include"`
	// 	} `json:"facets"`
	// 	SortBy string `json:"sortBy"`
	// 	Style  string `json:"style"`
	// } `json:"browse"`
	// ContentIsPost   bool        `json:"content_is_post"`
	// Description     string      `json:"description"`
	// ExpertResources interface{} `json:"expert_resources"`
	// FormFacets      struct {
	// } `json:"form_facets"`
	// ImageUrls   []string    `json:"image_urls"`
	// Next        interface{} `json:"next"`
	// NextSibling interface{} `json:"next_sibling"`
	// Pagination  struct {
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
	// Portal          bool        `json:"portal"`
	// Previous        interface{} `json:"previous"`
	// PreviousSibling interface{} `json:"previous_sibling"`
	Results []CollectionMetadata `json:"results"`
	// Search struct {
	// 	Dates       interface{} `json:"dates"`
	// 	FacetLimits string      `json:"facet_limits"`
	// 	Field       interface{} `json:"field"`
	// 	Hits        int         `json:"hits"`
	// 	In          string      `json:"in"`
	// 	Query       string      `json:"query"`
	// 	Recommended int         `json:"recommended"`
	// 	Site        struct {
	// 		Facets struct {
	// 			Include []struct {
	// 				Field string `json:"field"`
	// 				Label string `json:"label"`
	// 			} `json:"include"`
	// 		} `json:"facets"`
	// 		Mode   string `json:"mode"`
	// 		SortBy string `json:"sortBy"`
	// 		Style  string `json:"style"`
	// 	} `json:"site"`
	// 	SortBy      string `json:"sort_by"`
	// 	Type        string `json:"type"`
	// 	UnionFacets string `json:"union_facets"`
	// 	URL         string `json:"url"`
	// } `json:"search"`
	// Title string `json:"title"`
	// Type  string `json:"type"`
}

// CollectionMetadata describes a single collection
type CollectionMetadata struct {
	// AccessRestricted bool          `json:"access_restricted"`
	// Aka              []string      `json:"aka"`
	// Campaigns        []interface{} `json:"campaigns"`
	Contributor      []string  `json:"contributor"`
	Count            int       `json:"count"`
	Description      []string  `json:"description"`
	Digitized        bool      `json:"digitized"`
	ExtractTimestamp time.Time `json:"extract_timestamp"`
	Group            []string  `json:"group"`
	// Hassegments      bool          `json:"hassegments"`
	ID string `json:"id"`
	// ImageURL         []string      `json:"image_url"`
	// Index            int           `json:"index"`
	Item struct {
		AccessAdvisory []string `json:"access_advisory"`
		Contributors   []string `json:"contributors"`
		Date           string   `json:"date"`
		Format         []string `json:"format"`
		Language       []string `json:"language"`
		Location       []string `json:"location"`
		Medium         []string `json:"medium"`
		Notes          []string `json:"notes"`
		Repository     []string `json:"repository"`
		Subjects       []string `json:"subjects"`
		Summary        []string `json:"summary"`
		Title          string   `json:"title"`
	} `json:"item"`
	Items                string        `json:"items"`
	Language             []string      `json:"language"`
	Location             []string      `json:"location"`
	Number               []string      `json:"number"`
	NumberLccn           []string      `json:"number_lccn"`
	NumberSourceModified []string      `json:"number_source_modified"`
	OriginalFormat       []string      `json:"original_format"`
	OtherTitle           []interface{} `json:"other_title"`
	Partof               []string      `json:"partof"`
	Resources            []interface{} `json:"resources"`
	ShelfID              string        `json:"shelf_id"`
	Subject              []string      `json:"subject"`
	SubjectTopic         []string      `json:"subject_topic"`
	Timestamp            time.Time     `json:"timestamp"`
	Title                string        `json:"title"`
	URL                  string        `json:"url"`
	Dates                []time.Time   `json:"dates,omitempty"`
}

// String prints the title of the digital collection
func (cm CollectionMetadata) String() string {
	return cm.Title
}
