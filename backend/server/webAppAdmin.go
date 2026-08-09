package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
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

// checkGoogleSearchIndexingStatus checks if a URL is indexed in Google Search using the Search Console API.
// Returns true if indexed (verdict = "COVERED"), false otherwise.
func checkGoogleSearchIndexingStatus(ctx context.Context, url string) (*bool, error) {
	var credentialsPath = os.Getenv("GOOGLE_SEARCH_CONSOLE_CREDENTIALS")
	if credentialsPath == "" || !fileExists(credentialsPath) {
		return nil, fmt.Errorf("Google Search Console credentials not configured (GOOGLE_SEARCH_CONSOLE_CREDENTIALS env var must point to service account JSON)")
	}

	scope := "https://www.googleapis.com/auth/webmasters.readonly"

	// Use google.golang.org/api with a service account JWT token source.
	var searchService *searchconsole.Service
	var err error
	searchService, err = searchconsole.NewService(ctx,
		option.WithCredentialsFile(credentialsPath),
		option.WithScopes(scope))
	if err != nil {
		return nil, fmt.Errorf("failed to create Search Console service: %w", err)
	}

	var request = &searchconsole.InspectUrlIndexRequest{
		InspectionUrl: url,
		SiteUrl:       strings.TrimRight(url, "/"),
	}
	result, httpErr := searchService.UrlInspection.Index.Inspect(request).Context(ctx).Do()
	if httpErr != nil {
		return nil, fmt.Errorf("Search Console API error for %s: %w", url, httpErr)
	}

	var verdict = result.InspectionResult.IndexStatusResult.Verdict
	// COVERED means indexed; other values (NEUTRAL, NOT_ALLOWED, etc.) mean not indexed.
	var isIndexed = strings.EqualFold(verdict, "COVERED")
	return &isIndexed, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (me *webAppAdmin) getUrlPings(response http.ResponseWriter, request *http.Request) {
	var records = []rest_objects.GoalPostHeader{}
	var webLanguage = base.SupportedLanguages[0] // Only blog posts in main language currently participate in search indexing
	var publicBaseUrl = gophers.ReadEnvVar("PUBLIC_URL", default_public_url)
	var languagePath = webStaticGoals{}.getLanguagePath(webLanguage)
	me.db.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		if !row.SearchIndexingEnabled {
			return true
		}
		var record rest_objects.GoalPostHeader
		record.DateTime = row.DateTime
		record.IsPublic = row.IsPublic
		record.Type = row.TypeString
		record.Title = row.GetTranslatedTitle(webLanguage)
		record.Title = gophers.IfElse(record.Title != "", record.Title, row.TitleEnglish)
		record.GooglePingedAt = row.GooglePingedAt
		record.PublicUrl = publicBaseUrl + languagePath + "/personal-goals/" +
			gophers.GetStringFromInt64(row.GoalId) + "/" +
			gophers.GetStringFromInt64(row.DateTime) + ".html"

		// Check Google Search indexing status if credentials are available
		if os.Getenv("GOOGLE_SEARCH_CONSOLE_CREDENTIALS") != "" {
			isIndexed, err := checkGoogleSearchIndexingStatus(context.Background(), record.PublicUrl)
			if err != nil {
				// Log warning but continue processing other posts
				fmt.Printf("Warning: Failed to check indexing status for %s: %v\n", record.PublicUrl, err)
			} else if isIndexed != nil && *isIndexed {
				record.GoogleSearchIndexingStatus = isIndexed
			}
		}

		records = append(records, record)
		return true
	}, (db_objects.GoalPostRow{}).GetAllFieldSelector(), -1)
	writeJsonResponse(response, records)
}
