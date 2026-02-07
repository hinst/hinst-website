package server

import (
	"log"
	"net/http"
)

type webApp struct {
	db          *database
	allowOrigin string
	webPath     string
}

func (webApp) getDefaultWebPath() string {
	return "/hinst-website"
}

func (me *webApp) init(db *database) {
	me.db = db
	if me.webPath == "" {
		me.webPath = me.getDefaultWebPath()
	}
	if me.webPath == "/" {
		me.webPath = ""
	}
	var appGoals = new(webAppGoals)
	me.addFunctions(me.webPath, appGoals.init(me.db))

	var pageGoals = new(webPageGoals)
	me.addFunctions(me.webPath+pagesWebPath, pageGoals.init(me.db, me.webPath+pagesWebPath))

	var appAdmin = new(webAppAdmin)
	me.addFunctions(me.webPath, appAdmin.init(me.db))
}

func (me *webApp) addFunctions(path string, functions []namedWebFunction) {
	for _, namedWebFunction := range functions {
		var url = path + namedWebFunction.Name
		log.Printf("Adding web function: %v", url)
		http.HandleFunc(url, me.wrap(namedWebFunction.Function))
	}
}

func (me *webApp) wrap(function webFunction) webFunction {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", me.allowOrigin)
		defer func() {
			var exception = recover()
			if exception != nil {
				var webError, isWebError = exception.(webError)
				if isWebError {
					response.WriteHeader(webError.Status)
					var _, _ = response.Write(encodeJson(webError))
				} else {
					response.WriteHeader(http.StatusInternalServerError)
					panic(exception)
				}
			}
		}()
		function(response, request)
	}
}
