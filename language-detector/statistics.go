package main

import (
	"fmt"

	"github.com/jdkato/prose/v2"
)

type LanguageStats map[string]int

func (s LanguageStats) incrementKey(k string) {
	s[k] += 1
}

func tokenize(s string) (*prose.Document, error) {
	doc, err := prose.NewDocument(s,
		prose.WithExtraction(false),
		prose.WithTagging(false),
		prose.WithTokenization(false),
		prose.WithSegmentation(true))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// CalculateLanguages computes the number of sentences identified as each language.
// This is returned as a map with a count of sentences matching each language.
func CalculateLanguages(text string) (results LanguageStats, err error) {

	// Keep track of the number of sentences with each language
	results = make(LanguageStats)
	// whatlang = make(LanguageStats)

	doc, err := tokenize(text)
	if err != nil {
		return nil, fmt.Errorf("Error tokenizing text: %w", err)
	}

	// Count languages with both libraries
	for _, s := range doc.Sentences() {
		// wl := whatlanggo.Detect(s.Text)
		// if wl.IsReliable() {
		// 	whatlang.incrementKey(wl.Lang.Iso6393())
		// } else {
		// 	whatlang.incrementKey("und")
		// }

		lang, exists := ldetector.DetectLanguageOf(s.Text)
		if exists {
			results.incrementKey(lang.IsoCode639_3().String())
		} else {
			results.incrementKey("UND")
		}
	}

	return

}
