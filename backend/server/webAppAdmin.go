package server

import (
	"net/http"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"github.com/hinst/hinst-website/server/rest_objects"
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

func (me *webAppAdmin) getUrlPings(response http.ResponseWriter, request *http.Request) {
	var objects = []rest_objects.GoalPostSearchIndexingHeader{}
	var webLanguage = base.SupportedLanguages[0] // Only blog posts in main language currently participate in search indexing
	var languagePath = webStaticGoals{}.getLanguagePath(webLanguage)
	me.db.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		if !row.SearchIndexingEnabled {
			return true
		}
		var header rest_objects.GoalPostSearchIndexingHeader
		header.Read(row, webLanguage)
		header.PublicUrl = me.publicUrl() + languagePath + "/personal-goals/" +
			gophers.GetStringFromInt64(row.GoalId) + "/" +
			gophers.GetStringFromInt64(row.DateTime) + ".html"

		objects = append(objects, header)
		return true
	}, (db_objects.GoalPostRow{}).GetAllFieldSelector(), -1)
	writeJsonResponse(response, objects)
}

func (webAppAdmin) publicUrl() string {
	return gophers.ReadEnvVar("PUBLIC_URL", default_public_url)
}
