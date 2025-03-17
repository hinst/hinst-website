package main

import "net/http"

type webApp struct {
	savedGoalsPath string
	allowOrigin    string
	webPath        string
}

func (me *webApp) init() {
	if me.webPath == "" {
		me.webPath = "/hinst-website"
	}
	var goals = webAppGoals{savedGoalsPath: me.savedGoalsPath}
	for _, namedWebFunction := range goals.init() {
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
					panic(exception)
				}
			}
		}()
		function(response, request)
	}
}
