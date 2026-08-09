package server

import (
	"net/http"

	"github.com/hinst/hinst-website/server/db_objects"
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
	me.db.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		if !row.SearchIndexingEnabled {
			return true
		}
		return true
	}, (db_objects.GoalPostRow{}).GetAllFieldSelector(), -1)
	//TODO
}
