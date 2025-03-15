package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var supportedLanguages = []language.Tag{language.Russian, language.English, language.German}
var supportedLanguagesMatcher = language.NewMatcher(supportedLanguages)

func getLanguageName(tag language.Tag) string {
	return display.English.Languages().Name(tag)
}
