package database_objects

import (
	"slices"

	"github.com/hinst/go-common"
	"github.com/hinst/hinst-website/server/base"
	"golang.org/x/text/language"
)

func GetLanguagePostfix(supportedLanguage language.Tag) string {
	common.AssertCondition(slices.Contains(base.SupportedLanguages, supportedLanguage),
		func() string { return "Unsupported language: " + supportedLanguage.String() })
	var languageName = ""
	if supportedLanguage != base.SupportedLanguages[0] {
		languageName = base.GetLanguageName(supportedLanguage)
	}
	return languageName
}
