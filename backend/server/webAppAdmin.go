package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"github.com/hinst/hinst-website/server/rest_objects"
	"google.golang.org/api/option"
	"google.golang.org/api/searchconsole/v1"
)

type webAppAdmin struct {
	webAppBase
	db *database
}

func (me *webAppAdmin) init(db *database) []namedWebFunction {
	me.db = db
	var functions = []namedWebFunction{{"/api/urlPings", me.getUrlPings}}
	for i := range functions {
		functions[i].Function = me.guardAdminFunction(functions[i].Function)
	}
	return functions
}

func (me *webAppAdmin) checkSearchIndexing(ctx context.Context, url string) (string, error) {
	if me.googleAccountJson() == "" || !gophers.CheckFileExists(me.googleAccountJson()) {
		return "", fmt.Errorf("Need GOOGLE_ACCOUNT_JSON file: '%v'", me.googleAccountJson())
	}

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

func (me *webAppAdmin) getUrlPings(response http.ResponseWriter, request *http.Request) {
	var records = []rest_objects.GoalPostHeader{}
	var webLanguage = base.SupportedLanguages[0] // Only blog posts in main language currently participate in search indexing
	var languagePath = webStaticGoals{}.getLanguagePath(webLanguage)
	me.db.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		if !row.SearchIndexingEnabled {
			return true
		}
		var record rest_objects.GoalPostHeader
		record.DateTime = row.DateTime
		record.IsPublic = row.IsPublic
		record.Type = row.Type
		record.Title = row.GetTranslatedTitle(webLanguage)
		record.Title = gophers.IfElse(record.Title != "", record.Title, row.TitleEnglish)
		record.GooglePingedAt = row.GooglePingedAt
		record.PublicUrl = me.publicUrl() + languagePath + "/personal-goals/" +
			gophers.GetStringFromInt64(row.GoalId) + "/" +
			gophers.GetStringFromInt64(row.DateTime) + ".html"

		if me.googleAccountJson() != "" {
			googleSearchIndexingStatus, err := me.checkSearchIndexing(context.Background(), record.PublicUrl)
			if err != nil {
				log.Printf("Warning: Failed to check indexing status for %v: %v\n", record.PublicUrl, err)
			} else {
				record.GoogleSearchIndexingStatus = googleSearchIndexingStatus
			}
		}

		records = append(records, record)
		return true
	}, (db_objects.GoalPostRow{}).GetAllFieldSelector(), -1)
	writeJsonResponse(response, records)
}

// Path to JSON file containing authentication data for service account
func (webAppAdmin) googleAccountJson() string {
	return gophers.ReadEnvVar("GOOGLE_ACCOUNT_JSON", "")
}

func (webAppAdmin) publicUrl() string {
	return gophers.ReadEnvVar("PUBLIC_URL", default_public_url)
}

func (webAppAdmin) siteUrl() (siteUrl string) {
	siteUrl = webAppAdmin{}.publicUrl()
	if !strings.HasSuffix(siteUrl, "/") {
		siteUrl += "/"
	}
	return siteUrl
}
