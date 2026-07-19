package server

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/language"
)

type namedWebFunction struct {
	Name     string
	Function gophers.WebFunction
}

type WebRequest struct {
	Language      language.Tag
	WebPath       string
	StaticPath    string
	JpegExtension string
	HtmlExtension string
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
	var _, index, _ = base.SupportedLanguagesMatcher.Match([]language.Tag{tag}...)
	return base.SupportedLanguages[index]
}

func parseLanguageHeader(text string) language.Tag {
	var tags, _, parsedError = language.ParseAcceptLanguage(text)
	if parsedError != nil {
		panic(webError{"Invalid language header: " + text, http.StatusBadRequest})
	}
	gophers.AssertError(parsedError)
	var _, index, _ = base.SupportedLanguagesMatcher.Match(tags...)
	var tag = base.SupportedLanguages[index]
	return tag
}

func inputValidWebInteger(text string) int {
	var result, parseError = strconv.Atoi(text)
	var createWebError = func() webError {
		return webError{"Need integer. Received: " + text, http.StatusBadRequest}
	}
	gophers.AssertCondition(parseError == nil, createWebError)
	return result
}

func writeJsonResponse(response http.ResponseWriter, value any) {
	response.Header().Set(gophers.ContentTypeHeader, gophers.ContentTypeJson)
	var _, _ = response.Write(gophers.EncodeJson(value))
}

func writeHtmlResponse(response http.ResponseWriter, text string) {
	text = gophers.AssertResultError(formatHtml(text))
	response.Header().Set(gophers.ContentTypeHeader, "text/html; charset=utf-8")
	var _, _ = response.Write([]byte(text))
}

func readTextFromUrl(url string) string {
	return string(readBytesFromUrl(url))
}

func readBytesFromUrl(url string) []byte {
	var response = gophers.AssertResultError(http.Get(url))
	defer gophers.IoCloseSilently(response.Body)
	assertResponse(response)
	var data = gophers.AssertResultError(io.ReadAll(response.Body))
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

// Format HTML using Prettier server.
// Returns error if supplied HTML is invalid.
// Panics if unable to connect to Prettier server.
func formatHtml(text string) (string, error) {
	var client = &http.Client{Timeout: 10 * time.Minute}
	var url = gophers.RequireEnvVar("PRETTIER_SERVER_URL") +
		gophers.BuildUrlQueryParams(map[string]string{"filename": "index.html"})
	var textBytes = []byte(text)
	var request = gophers.AssertResultError(http.NewRequest("POST", url, bytes.NewBuffer(textBytes)))
	request.Header.Set(gophers.ContentTypeHeader, "text/html")
	var response = gophers.AssertResultError(client.Do(request))
	defer gophers.IoCloseSilently(response.Body)
	var responseBytes = gophers.AssertResultError(io.ReadAll(response.Body))
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

func getUrlBase64(contentType string, array []byte) string {
	return "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(array)
}

func stripHtml(text string) string {
	return bluemonday.StrictPolicy().Sanitize(text)
}
