package main

import (
	"github.com/abadojack/whatlanggo"
	"github.com/jdkato/prose/v2"
)

type LanguageStats map[string]int

func CalculateLanguages(s string) (lingua LanguageStats, whatlang LanguageStats, err error) {

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

	for _, s := range doc.Sentences() {
		wl := whatlanggo.Detect(s.Text)
		if wl.IsReliable() {
			whatlang[wl.Lang.String()] = whatlang[wl.Lang.String()] + 1
		} else {
			whatlang["Unknown"] = whatlang["Unknown"] + 1
		}
	}

	return

}
