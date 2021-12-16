package main

import "github.com/pemistahl/lingua-go"

var ldetector = lingua.NewLanguageDetectorBuilder().FromAllLanguages().Build()
