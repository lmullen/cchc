package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguageStats_incrementKey(t *testing.T) {
	t.Parallel()
	stats := make(LanguageStats)
	k := "test"
	for i := 0; i < 100; i++ {
		stats.incrementKey(k)
	}
	assert.Equal(t, 100, stats[k])
}

func TestCalculateLanguages(t *testing.T) {
	t.Parallel()

	text := "This is a dummy document. Ese reloj fue un regalo de mi mujer.  It has two English, one Spanish, and one German sentences. Ich möchte ein Bier."

	results := make(LanguageStats)

	err := CalculateLanguages(text, results)
	assert.NoError(t, err)

	expected := LanguageStats{
		"DEU": 1,
		"ENG": 2,
		"SPA": 1,
	}

	assert.Equal(t, expected, results)
}
