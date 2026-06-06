package base

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var SupportedLanguages = []language.Tag{language.Russian, language.English, language.German}
var SupportedLanguagesMatcher = language.NewMatcher(SupportedLanguages)

func GetLanguageName(tag language.Tag) string {
	return display.English.Languages().Name(tag)
}
