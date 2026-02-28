package server

import (
	"bytes"
	"context"
	"net/http"

	"github.com/hinst/go-common"
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
	var jsonText = common.ReadBytesFile(requireEnvVar("GOOGLE_ACCOUNT_JSON"))
	var conf = common.AssertResultError(google.JWTConfigFromJSON(jsonText, me.getScope()))
	me.client = conf.Client(context.Background())
}

func (me *GoogleIndexingClient) updateUrl(url string) bool {
	var data = GoogleUrlNotification{
		Url:  url,
		Type: "URL_UPDATED",
	}
	var apiUrl = "https://indexing.googleapis.com/v3/urlNotifications:publish"
	var response = common.AssertResultError(me.client.Post(apiUrl,
		common.ContentTypeJson, bytes.NewReader(common.EncodeJson(data))))
	defer common.IoCloseSilently(response.Body)
	if response.StatusCode == http.StatusTooManyRequests {
		return false
	}
	assertResponse(response)
	return true
}
