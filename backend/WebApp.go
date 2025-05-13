package main

import "net/http"

type webApp struct {
	db             *database
	savedGoalsPath string
	allowOrigin    string
	webPath        string
}

func (me *webApp) init(db *database) {
	me.db = db
	if me.webPath == "" {
		me.webPath = "/hinst-website"
	}
	var appGoals = new(webAppGoals)
	appGoals.savedGoalsPath = me.savedGoalsPath
	me.addFunctions(appGoals.init(me.db))
	var appRiddles = new(webAppRiddles)
	me.addFunctions(appRiddles.init(me.db))
}

func (me *webApp) addFunctions(functions []namedWebFunction) {
	for _, namedWebFunction := range functions {
		http.HandleFunc(me.webPath+namedWebFunction.Name, me.wrap(namedWebFunction.Function))
	}
}

func (me *webApp) wrap(function func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", me.allowOrigin)
		defer func() {
			var exception = recover()
			if exception != nil {
				var webError, isWebError = exception.(webError)
				if isWebError {
					response.WriteHeader(webError.Status)
					response.Write(encodeJson(webError))
				} else {
					response.WriteHeader(http.StatusInternalServerError)
					panic(exception)
				}
			}
		}()
		function(response, request)
	}
}
