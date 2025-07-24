package main

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/goccy/go-json"
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

func requireRequestQueryInt(request *http.Request, key string) int {
	var text = request.URL.Query().Get(key)
	if text == "" {
		panic(webError{Message: "Missing integer parameter '" + key + "'", Status: http.StatusBadRequest})
	}
	var value, parseError = strconv.Atoi(text)
	assertCondition(parseError == nil, func() error {
		return webError{Message: "Invalid integer: " + text, Status: http.StatusBadRequest}
	})
	return value
}

func decodeWebJson(input io.ReadCloser, value any) {
	var decodeError = json.NewDecoder(input).Decode(value)
	if decodeError != nil {
		panic(webError{Message: "Invalid JSON body: " + decodeError.Error(), Status: http.StatusBadRequest})
	}
}

func writeJsonResponse(response http.ResponseWriter, value any) {
	response.Header().Set("Content-Type", contentTypeJson)
	response.Write(encodeJson(value))
}
