package items

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ItemResponse represents an item-level object returned from the API. Many more fields
// are returned and will be stored in the database as a JSONB field, but these
// are the ones that will be serialized to regular database fields.
type ItemResponse struct {
	ItemDetails struct {
		ID       string   `json:"id"`
		URL      string   `json:"url"`
		Date     string   `json:"date"`
		Subjects []string `json:"subject_headings"`
		Title    string   `json:"title"`
		Language []string `json:"language"`
		// OnlineFormat []string  `json:"online_format"`
		// Version      int64     `json:"_version_"`
		// HasSegments  bool      `json:"hassegments"`
	} `json:"item"`
	Resources []struct {
		FulltextFile string `json:"fulltext_file,omitempty"`
		DJVUTextFile string `json:"djvu_text_file,omitempty"`
		Image        string `json:"image,omitempty"`
		PDF          string `json:"pdf,omitempty"`
		URL          string `json:"url,omitempty"`
		Caption      string `json:"caption,omitempty"`
		Files        [][]struct {
			Mimetype        string `json:"mimetype,omitempty"`
			Fulltext        string `json:"fulltext,omitempty"`
			FulltextService string `json:"fulltext_service,omitempty"`
			WordCoordinates string `json:"word_coordinates,omitempty"`
			URL             string `json:"url,omitempty"`
			Info            string `json:"info,omitempty"`
			Use             string `json:"use,omitempty"`
		} `json:"files"`
	} `json:"resources"`
}

// Fetch gets an item's metadata from the LOC.gov API.
func (i *Item) Fetch(client *http.Client) error {

	u, _ := url.Parse(i.URL.String)
	remove := []string{"more_like_this", "related_items", "cite_this", "options"}
	options := url.Values{
		"at!": []string{strings.Join(remove, ",")},
		"fo":  []string{"json"},
	}
	u.RawQuery = options.Encode()
	url := u.String()

	response, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("Error getting item over HTTP: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error reading HTTP response body while fetching item: %w", err)
	}

	var result ItemResponse

	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("Error unmarshalling item metadata: %w", err)
	}

	i.ID = result.ItemDetails.ID
	i.URL.Scan(result.ItemDetails.URL)
	i.Title.Scan(result.ItemDetails.Title)
	i.Year = year(result.ItemDetails.Date)
	i.Date.Scan(result.ItemDetails.Date)
	i.Subjects = result.ItemDetails.Subjects
	i.Languages = result.ItemDetails.Language
	i.API.Scan(data)

	// Iterate through all the files and formats to get the full text representations
	for resourceSeq, resource := range result.Resources {
		var r ItemResource
		r.ItemID = i.ID
		r.ResourceSeq = resourceSeq
		if resource.FulltextFile != "" {
			r.FullTextFile.Scan(resource.FulltextFile)
		}
		if resource.DJVUTextFile != "" {
			r.DJVUTextFile.Scan(resource.DJVUTextFile)
		}
		if resource.Image != "" {
			r.Image.Scan(resource.Image)
		}
		if resource.PDF != "" {
			r.PDF.Scan(resource.PDF)
		}
		if resource.URL != "" {
			r.URL.Scan(resource.URL)
		}
		if resource.Caption != "" {
			r.Caption.Scan(resource.Caption)
		}
		i.Resources = append(i.Resources, r)
		for fileSeq, file := range resource.Files {
			for formatSeq, format := range file {
				var f ItemFile
				f.ItemID = i.ID
				f.ResourceSeq = resourceSeq
				f.FileSeq = fileSeq
				f.FormatSeq = formatSeq
				if format.Mimetype != "" {
					f.Mimetype.Scan(format.Mimetype)
				}
				if format.Fulltext != "" {
					f.FullText.Scan(format.Fulltext)
				}
				if format.FulltextService != "" {
					f.FullTextService.Scan(format.FulltextService)
				}
				if format.WordCoordinates != "" {
					f.WordCoordinates.Scan(format.WordCoordinates)
				}
				if format.URL != "" {
					f.URL.Scan(format.URL)
				}
				if format.Info != "" {
					f.Info.Scan(format.Info)
				}
				if format.Use != "" {
					f.Use.Scan(format.Use)
				}
				i.Files = append(i.Files, f)
			}
		}
	}

	return nil

}
