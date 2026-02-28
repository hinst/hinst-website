package server

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hinst/go-common"
	"golang.org/x/text/language"
)

type namedWebFunction struct {
	Name     string
	Function common.WebFunction
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
	common.AssertError(parsedError)
	var _, index, _ = supportedLanguagesMatcher.Match(tags...)
	var tag = supportedLanguages[index]
	return tag
}

func setCacheAge(response http.ResponseWriter, duration time.Duration) {
	response.Header().Set(common.CacheControlHeader, "max-age="+strconv.Itoa(int(duration.Seconds())))
}

func inputValidWebInteger(text string) int {
	var result, parseError = strconv.Atoi(text)
	var createWebError = func() webError {
		return webError{"Need integer. Received: " + text, http.StatusBadRequest}
	}
	common.AssertCondition(parseError == nil, createWebError)
	return result
}

func writeJsonResponse(response http.ResponseWriter, value any) {
	response.Header().Set(common.ContentTypeHeader, common.ContentTypeJson)
	var _, _ = response.Write(common.EncodeJson(value))
}

func writeHtmlResponse(response http.ResponseWriter, text string) {
	text = common.AssertResultError(formatHtml(text))
	response.Header().Set(common.ContentTypeHeader, "text/html; charset=utf-8")
	var _, _ = response.Write([]byte(text))
}

func readTextFromUrl(url string) string {
	return string(readBytesFromUrl(url))
}

func readBytesFromUrl(url string) []byte {
	var response = common.AssertResultError(http.Get(url))
	defer ioCloseSilently(response.Body)
	assertResponse(response)
	var data = common.AssertResultError(io.ReadAll(response.Body))
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

func formatHtml(text string) (string, error) {
	var client = &http.Client{Timeout: 10 * time.Minute}
	var url = requireEnvVar("PRETTIER_SERVER_URL") +
		buildUrl("", map[string]string{"filename": "index.html"})
	var textBytes = []byte(text)
	var request = common.AssertResultError(http.NewRequest("POST", url, bytes.NewBuffer(textBytes)))
	request.Header.Set(common.ContentTypeHeader, "text/html")
	var response = common.AssertResultError(client.Do(request))
	defer ioCloseSilently(response.Body)
	var responseBytes = common.AssertResultError(io.ReadAll(response.Body))
	var responseText = string(responseBytes)
	if response.StatusCode == http.StatusOK {
		return responseText, nil
	} else {
		var errorText = "Cannot format HTML text; status: " + response.Status + "; response: " + responseText
		if response.StatusCode == http.StatusBadRequest {
			// The supplied text was not a valid HTML
			return text, errors.New(errorText)
		} else {
			panic(errorText)
		}
	}
}

func validateHtml(text string) error {
	var _, e = formatHtml("<div>" + text + "</div>")
	return e
}
