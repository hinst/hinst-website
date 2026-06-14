package server

import (
	"bytes"
	"context"
	"net/http"

	"github.com/hinst/go-gophers"
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
	var jsonText = gophers.ReadBytesFile(gophers.RequireEnvVar("GOOGLE_ACCOUNT_JSON"))
	var conf = gophers.AssertResultError(google.JWTConfigFromJSON(jsonText, me.getScope()))
	me.client = conf.Client(context.Background())
}

func (me *GoogleIndexingClient) updateUrl(url string) bool {
	var data = GoogleUrlNotification{
		Url:  url,
		Type: "URL_UPDATED",
	}
	var apiUrl = "https://indexing.googleapis.com/v3/urlNotifications:publish"
	var response = gophers.AssertResultError(me.client.Post(apiUrl,
		gophers.ContentTypeJson, bytes.NewReader(gophers.EncodeJson(data))))
	defer gophers.IoCloseSilently(response.Body)
	if response.StatusCode == http.StatusTooManyRequests {
		return false
	}
	assertResponse(response)
	return true
}
