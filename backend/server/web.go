package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/yosssi/gohtml"
	"golang.org/x/text/language"
)

const contentTypeJson = "application/json"
const contentTypeText = "text/plain"

type webFunction func(response http.ResponseWriter, request *http.Request)
type namedWebFunction struct {
	Name     string
	Function webFunction
}

func getWebLanguage(request *http.Request) language.Tag {
	var queryLanguage = request.URL.Query().Get("lang")
	if len(queryLanguage) > 0 {
		return parseLanguageTag(queryLanguage)
	}
	var acceptLanguage = request.Header.Get("Accept-Language")
	return parseLanguageHeader(acceptLanguage)
}

func parseLanguageTag(text string) language.Tag {
	var tag, parsedError = language.Parse(text)
	if parsedError != nil {
		panic(webError{"Invalid language tag: " + text, http.StatusBadRequest})
	}
	var _, index, _ = supportedLanguagesMatcher.Match([]language.Tag{tag}...)
	return supportedLanguages[index]
}

func parseLanguageHeader(text string) language.Tag {
	var tags, _, parsedError = language.ParseAcceptLanguage(text)
	if parsedError != nil {
		panic(webError{"Invalid language header: " + text, http.StatusBadRequest})
	}
	assertError(parsedError)
	var _, index, _ = supportedLanguagesMatcher.Match(tags...)
	var tag = supportedLanguages[index]
	return tag
}

func setCacheAge(response http.ResponseWriter, duration time.Duration) {
	response.Header().Set("Cache-Control", "max-age="+strconv.Itoa(int(duration.Seconds())))
}

func inputValidWebInteger(text string) int {
	var result, parseError = strconv.Atoi(text)
	var createWebError = func() webError {
		return webError{"Need integer. Received: " + text, http.StatusBadRequest}
	}
	assertCondition(parseError == nil, createWebError)
	return result
}

func decodeWebJson(input io.ReadCloser, value any) {
	var decodeError = json.NewDecoder(input).Decode(value)
	if decodeError != nil {
		panic(webError{Message: "Invalid JSON body: " + decodeError.Error(), Status: http.StatusBadRequest})
	}
}

func writeJsonResponse(response http.ResponseWriter, value any) {
	response.Header().Set("Content-Type", contentTypeJson)
	var _, _ = response.Write(encodeJson(value))
}

func writeHtmlResponse(response http.ResponseWriter, text string) {
	text = gohtml.Format(text)
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	var _, _ = response.Write([]byte(text))
}

func readTextFromUrl(url string) string {
	return string(readBytesFromUrl(url))
}

func readBytesFromUrl(url string) []byte {
	var response = assertResultError(http.Get(url))
	defer ioCloseSilently(response.Body)
	assertResponse(response)
	var data = assertResultError(io.ReadAll(response.Body))
	return data
}

func assertResponse(response *http.Response) {
	if response.StatusCode != http.StatusOK {
		var text, _ = io.ReadAll(response.Body)
		panic("Bad status=" + response.Status +
			" returned from url=" + response.Request.URL.String() +
			"\n" + string(text))
	}
}

func buildUrl(base string, parameters map[string]string) string {
	var theUrl = base
	var first = true
	for key, value := range parameters {
		if first {
			theUrl += "?"
			first = false
		} else {
			theUrl += "&"
		}
		theUrl += key + "=" + url.QueryEscape(value)
	}
	return theUrl
}
