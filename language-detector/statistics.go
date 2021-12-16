package main

import (
	"github.com/abadojack/whatlanggo"
	"github.com/jdkato/prose/v2"
)

type LanguageStats map[string]int

func CalculateLanguages(s string) (lingua LanguageStats, whatlang LanguageStats, err error) {

	// Keep track of the number of sentences with each language
	lingua = make(LanguageStats)
	whatlang = make(LanguageStats)

	doc, err := prose.NewDocument(s,
		prose.WithExtraction(false),
		prose.WithTagging(false),
		prose.WithTokenization(false),
		prose.WithSegmentation(true))
	if err != nil {
		return nil, nil, err
	}

	// Count languages using whatlang
	for _, s := range doc.Sentences() {
		wl := whatlanggo.Detect(s.Text)
		if wl.IsReliable() {
			whatlang[wl.Lang.String()] = whatlang[wl.Lang.String()] + 1
		} else {
			whatlang["Unknown"] = whatlang["Unknown"] + 1
		}
	}

	// Count languages using lingua
	for _, s := range doc.Sentences() {
		language, exists := ldetector.DetectLanguageOf(s.Text)
		if exists {
			lingua[language.String()] = lingua[language.String()] + 1
		} else {
			lingua["Unknown"] = lingua["Unknown"] + 1
		}
	}

	return

}
