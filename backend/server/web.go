package server

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/language"
)

const default_public_url = "https://hinst.github.io"

var _webRouter *http.ServeMux

func webRouter() *http.ServeMux {
	if _webRouter == nil {
		_webRouter = http.NewServeMux()
	}
	return _webRouter
}

type namedWebFunction struct {
	Name     string
	Function gophers.WebFunction
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

func writeJsonResponse(response http.ResponseWriter, value any) {
	response.Header().Set(gophers.ContentTypeHeader, gophers.ContentTypeJson)
	var _, _ = response.Write(gophers.EncodeJson(value))
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
