package main

import "time"

func checkable(failures map[string]time.Time, key string) bool {
	timeFailed, previouslyFailed := failures[key]

	// If it hasn't previously failed, then it is new and is checkable
	if !previouslyFailed {
		return true
	}

	// Failed more than an hour ago, so it is checkable
	if time.Now().Sub(timeFailed).Hours() >= 1.0 {
		return true
	}

	// Otherwise not checkable
	return false

}
