package main

import (
	"net/http"
	"strconv"
	"time"

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

func setCacheAge(response http.ResponseWriter, duration time.Duration) {
	response.Header().Set("Cache-Control", "max-age="+strconv.Itoa(int(duration.Seconds())))
}
