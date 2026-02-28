package server

import (
	"bytes"
	"context"
	"net/http"

	"golang.org/x/oauth2/google"
)

type GoogleIndexingClient struct {
	client *http.Client
}

type GoogleUrlNotification struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}

func (GoogleIndexingClient) getScope() string {
	return "https://www.googleapis.com/auth/indexing"
}

func (me *GoogleIndexingClient) connect() {
	var jsonText = readBytesFile(requireEnvVar("GOOGLE_ACCOUNT_JSON"))
	var conf = AssertResultError(google.JWTConfigFromJSON(jsonText, me.getScope()))
	me.client = conf.Client(context.Background())
}

func (me *GoogleIndexingClient) updateUrl(url string) bool {
	var data = GoogleUrlNotification{
		Url:  url,
		Type: "URL_UPDATED",
	}
	var apiUrl = "https://indexing.googleapis.com/v3/urlNotifications:publish"
	var response = AssertResultError(me.client.Post(apiUrl,
		contentTypeJson, bytes.NewReader(encodeJson(data))))
	defer ioCloseSilently(response.Body)
	if response.StatusCode == http.StatusTooManyRequests {
		return false
	}
	assertResponse(response)
	return true
}
