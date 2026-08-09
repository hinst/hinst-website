package server

import (
	"net/http"

	"github.com/hinst/go-gophers"
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
	var records = make([]rest_objects.GoalPostHeader, 0)
	me.db.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		if !row.SearchIndexingEnabled {
			return true
		}
		var record rest_objects.GoalPostHeader
		record.DateTime = row.DateTime
		record.IsPublic = row.IsPublic
		record.Type = row.TypeString
		record.Title = row.GetTranslatedTitle(getWebLanguage(request))
		record.Title = gophers.IfElse(record.Title != "", record.Title, row.TitleEnglish)
		record.GooglePingedAt = row.GooglePingedAt
		records = append(records, record)
		return true
	}, (db_objects.GoalPostRow{}).GetAllFieldSelector(), -1)
	writeJsonResponse(response, records)
}
