package base

import (
	"slices"

	"github.com/hinst/go-gophers"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

const JOURNALS_TITLE = "Journals"

var SupportedLanguages = []language.Tag{language.Russian, language.English, language.German}
var SupportedLanguagesMatcher = language.NewMatcher(SupportedLanguages)

func GetLanguageName(tag language.Tag) string {
	return display.English.Languages().Name(tag)
}

func GetLanguageIndex(tag language.Tag) int {
	var index = slices.Index(SupportedLanguages, tag)
	gophers.AssertCondition(index >= 0, func() string {
		return "Unsupported language: " + tag.String()
	})
	return index
}
