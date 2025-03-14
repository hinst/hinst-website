package main

import (
	"net/http"

	"golang.org/x/text/language"
)

const contentTypeJson = "application/json"

type webFunction func(response http.ResponseWriter, request *http.Request)
type namedWebFunction struct {
	Name     string
	Function webFunction
}

func getWebLanguage(request *http.Request) language.Tag {
	var acceptLanguage = request.Header.Get("Accept-Language")
	var tags, _, parsedError = language.ParseAcceptLanguage(acceptLanguage)
	assertError(parsedError)
	var _, index, _ = supportedLanguagesMatcher.Match(tags...)
	var tag = supportedLanguages[index]
	return tag
}
