package jobs

import "errors"

// ErrAllQueued is returned when there is need to enqueue any more items for a
// destination.
var ErrAllQueued = errors.New("All items have been queued for this destination")

var ErrNoJobs = errors.New("There are no ready jobs for that destination")
