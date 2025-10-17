package server

import "net/http"

type webAppAdmin struct {
	webAppBase
	db *database
}

func (me *webAppAdmin) init(db *database) []namedWebFunction {
	me.db = db
	return []namedWebFunction{
		{"/api/urlPings", me.guardAdminFunction(me.getUrlPings)},
	}
}

func (me *webAppAdmin) getUrlPings(response http.ResponseWriter, request *http.Request) {
}
