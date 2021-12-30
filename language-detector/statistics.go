package main

import (
	"fmt"

	"github.com/jdkato/prose/v2"
)

// LanguageStats is a map with an ISO 693-3 code for langauges and the count of
// the number of sentences in that language. Unknown languages are recorded as
// `UND`.
type LanguageStats map[string]int

func (s LanguageStats) incrementKey(k string) {
	s[k]++
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
// You pass in a map of the the languages. This is passed in rather than generated
// in the function, because we might want to count multiple pages in the same item.
func CalculateLanguages(text string, results LanguageStats) error {

	doc, err := tokenize(text)
	if err != nil {
		return fmt.Errorf("Error tokenizing text: %w", err)
	}

	// Detect language for each sentence and track results
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

	return nil

}
