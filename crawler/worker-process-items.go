package main

// StartProcessingItems begins reading items from the message queue in order to
// process each of them.
func StartProcessingItems() {
	for msg := range app.ItemMetadataQ.Consumer {
		// Give each item its own goroutine
		ProcessItemMetadata(msg)
	}
}
