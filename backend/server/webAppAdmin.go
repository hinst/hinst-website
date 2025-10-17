package server

import (
	"net/http"
	"time"
)

type webAppAdmin struct {
	webAppBase
	db *database
}

func (me *webAppAdmin) init(db *database) []namedWebFunction {
	me.db = db
	var functions = []namedWebFunction{
		{"/api/urlPings", me.getUrlPings},
		{"/api/pingUrlManually", me.pingUrlManually},
	}
	for i := range functions {
		functions[i].Function = me.guardAdminFunction(functions[i].Function)
	}
	return functions
}

func (me *webAppAdmin) getUrlPings(response http.ResponseWriter, request *http.Request) {
	var records = me.db.getUrlPings()
	writeJsonResponse(response, records)
}

func (me *webAppAdmin) pingUrlManually(response http.ResponseWriter, request *http.Request) {
	var url = request.URL.Query().Get("url")
	if url == "" {
		http.Error(response, "Missing URL parameter", http.StatusBadRequest)
		return
	}
	var success = me.db.updateUrlPingGoogleManually(url, time.Now())
	if !success {
		http.Error(response, "URL not found", http.StatusNotFound)
		return
	}
	writeJsonResponse(response, true)
}
