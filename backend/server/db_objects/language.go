package db_objects

import (
	"github.com/hinst/hinst-website/server/base"
	"golang.org/x/text/language"
)

// Returns a slice aligned with base.SupportedLanguages, where value is stored at the given language position
func NewLocalizedStringSlice(value string, languageTag language.Tag) []string {
	var result = make([]string, len(base.SupportedLanguages))
	result[base.GetLanguageIndex(languageTag)] = value
	return result
}
