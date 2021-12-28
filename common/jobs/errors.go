package jobs

import "errors"

// ErrAllQueued is returned when there is need to enqueue any more items for a
// destination. This would be an expected error.
var ErrAllQueued = errors.New("All items have been queued for this destination")

// ErrNoJobs is returned when there are no ready jobs for a particular destination.
// This would be an expected error.
var ErrNoJobs = errors.New("There are no ready jobs for that destination")
