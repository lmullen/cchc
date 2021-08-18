package main

import "strings"

// isResourceNotItem checks whether an item's URL is actually to a resource.
func isResourceNotItem(url string) bool {
	return strings.Contains(url, "/resource/")
}
