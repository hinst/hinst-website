package server

import (
	"context"

	"golang.org/x/oauth2/google"
)

const GOOGLE_SCOPE_INDEXING = "https://www.googleapis.com/auth/indexing"

func connectGoogleIndexing() {
	var jsonText = readBytesFile(requireEnvVar("GOOGLE_ACCOUNT_JSON"))
	var conf = assertResultError(google.ConfigFromJSON(jsonText, GOOGLE_SCOPE_INDEXING))
	client := conf.Client(context.Background(), nil)
	client.Get("...")
}
