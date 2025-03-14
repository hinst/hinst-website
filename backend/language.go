package main

import "golang.org/x/text/language"

var supportedLanguages = []language.Tag{language.Russian, language.English, language.German}
var supportedLanguagesMatcher = language.NewMatcher(supportedLanguages)
