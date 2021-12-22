package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_checkable(t *testing.T) {
	failures := map[string]time.Time{
		"failed-two-hours-ago":       time.Now().Add(-2 * time.Hour),
		"failed-fifteen-minutes-ago": time.Now().Add(-15 * time.Minute),
	}

	assert.True(t, checkable(failures, "failed-two-hours-ago"), "if it failed more than an hour ago check it")
	assert.False(t, checkable(failures, "failed-fifteen-minutes-ago"), "if it failed less than an hour ago, don't check it")
	assert.True(t, checkable(failures, "never-failed"), "if it never failed, check it")
}
