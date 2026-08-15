package server

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"google.golang.org/api/option"
	"google.golang.org/api/searchconsole/v1"
)

type searchIndexingUpdater struct {
	db *database
}

// Checks Google Search indexing status for all goal posts with search indexing enabled
// and saves the result together with the check timestamp into the database.
func (me *searchIndexingUpdater) run() {
	if me.googleAccountJson() == "" || !gophers.CheckFileExists(me.googleAccountJson()) {
		log.Printf("Search indexing status update skipped: need GOOGLE_ACCOUNT_JSON file: '%v'\n",
			me.googleAccountJson())
		return
	}
	var webLanguage = base.SupportedLanguages[0] // Only blog posts in main language currently participate in search indexing
	me.db.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		if !row.SearchIndexingEnabled {
			return true
		}
		var postPublicUrl = webStaticGoals{}.getPublicUrl(row, webLanguage)
		var searchIndexingStatus, searchIndexingError = me.checkSearchIndexing(context.Background(), postPublicUrl)
		if searchIndexingError != nil {
			log.Printf("Warning: Failed to check indexing status for %v: %v\n", postPublicUrl, searchIndexingError)
			return true
		}
		me.db.setGoalPostSearchIndexingStatus(
			row.GoalId, row.GetDateTime(), searchIndexingStatus, time.Now().UTC())
		return true
	}, (db_objects.GoalPostRow{}).GetAllFieldSelector(), -1)
}

func (me *searchIndexingUpdater) checkSearchIndexing(ctx context.Context, url string) (string, error) {
	searchService, searchConsoleError := searchconsole.NewService(ctx,
		option.WithAuthCredentialsFile(option.ServiceAccount, me.googleAccountJson()),
		option.WithScopes("https://www.googleapis.com/auth/webmasters.readonly"))
	if searchConsoleError != nil {
		return "", searchConsoleError
	}
	var request = &searchconsole.InspectUrlIndexRequest{
		InspectionUrl: url,
		SiteUrl:       me.siteUrl(),
	}
	result, searchConsoleError := searchService.UrlInspection.Index.Inspect(request).Context(ctx).Do()
	if searchConsoleError != nil {
		return "", searchConsoleError
	}
	return result.InspectionResult.IndexStatusResult.Verdict, nil
}

// Path to JSON file containing authentication data for Google service account
// The e-mail of the account should be added to Google Search Console -> Settings -> Users and permissions
func (searchIndexingUpdater) googleAccountJson() string {
	return gophers.ReadEnvVar("GOOGLE_ACCOUNT_JSON", "")
}

func (searchIndexingUpdater) publicUrl() string {
	return gophers.ReadEnvVar("PUBLIC_URL", default_public_url)
}

func (me searchIndexingUpdater) siteUrl() (siteUrl string) {
	siteUrl = me.publicUrl()
	if !strings.HasSuffix(siteUrl, "/") {
		siteUrl += "/"
	}
	return siteUrl
}
