package items

import "fmt"

// String
func (item Item) String() string {
	return fmt.Sprintf("[Item: %s. Resources: %v. Files: %v.]", item.ID, len(item.Resources), len(item.Files))
}

// Fetched reports whether the full metadata has been fetched for an item.
func (item *Item) Fetched() bool {
	return item.API.Valid
}
