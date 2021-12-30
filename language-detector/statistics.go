package main

import (
	"fmt"

	"github.com/abadojack/whatlanggo"
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

// CalculateLanguages computes the number of sentences identified as each language
// using two different libraries. These are returned as maps of the languages and the
// count of sentences.
func CalculateLanguages(text string) (lingua LanguageStats, whatlang LanguageStats, err error) {

	// Keep track of the number of sentences with each language
	lingua = make(LanguageStats)
	whatlang = make(LanguageStats)

	doc, err := tokenize(text)
	if err != nil {
		return nil, nil, fmt.Errorf("Error tokenizing text: %w", err)
	}

	// Count languages with both libraries
	for _, s := range doc.Sentences() {
		wl := whatlanggo.Detect(s.Text)
		if wl.IsReliable() {
			whatlang.incrementKey(wl.Lang.Iso6393())
		} else {
			whatlang.incrementKey("und")
		}

		lang, exists := ldetector.DetectLanguageOf(s.Text)
		if exists {
			lingua.incrementKey(lang.IsoCode639_3().String())
		} else {
			lingua.incrementKey("und")
		}

	}

	return

}
