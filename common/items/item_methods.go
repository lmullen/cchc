package items

import (
	"fmt"

	"github.com/k3a/html2text"
)

// String
func (item Item) String() string {
	return fmt.Sprintf("[Item: %s. Resources: %v. Files: %v.]", item.ID, len(item.Resources), len(item.Files))
}

// Fetched reports whether the full metadata has been fetched for an item.
func (item *Item) Fetched() bool {
	return item.API.Valid
}

// FullText returns a slice of plain text objects containing the page-level
// text, cleaned up to be plain text. The boolean return value keeps track of
// whether there is or isn't text available at this level.
func (item *Item) FullText() (text []PlainText, has bool) {

	// We have to look for plain text in a lot of places, so keep track of whether
	// we have found it yet.
	has = false

	// FULL TEXT CHECK 1: Has text/plain mimetype with fulltext field
	if !has {
		for _, file := range item.Files {
			if file.Mimetype.Valid && file.Mimetype.String == "text/plain" && file.FullText.Valid {
				text = append(text, PlainText{Text: file.FullText.String})
				has = true // Keep track that we have found full text
			}
		}
	}

	// FULL TEXT CHECK 2: Has text/xml mimetype with fulltext field
	if !has {
		for _, file := range item.Files {
			if file.Mimetype.Valid && file.Mimetype.String == "text/xml" && file.FullText.Valid {
				text = append(text, PlainText{Text: html2text.HTML2Text(file.FullText.String)})
				// text = append(text, PlainText{Text: stripXML.Sanitize(file.FullText.String)})
				has = true // Keep track that we have found full text
			}
		}

	}

	return text, has
}
