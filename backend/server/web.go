package server

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/yosssi/gohtml"
	"golang.org/x/text/language"
)

const contentTypeJson = "application/json"

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

func writeHtmlResponse(response http.ResponseWriter, text string) {
	text = gohtml.Format(text)
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.Write([]byte(text))
}

func readTextFromUrl(url string) string {
	var response = assertResultError(http.Get(url))
	if response.StatusCode != http.StatusOK {
		panic("Bad status=" + response.Status + " returned from url=" + url)
	}
	defer response.Body.Close()
	var text = string(assertResultError(io.ReadAll(response.Body)))
	return text
}
