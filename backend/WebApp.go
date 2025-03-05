package main

import (
	"net/http"
	"os"
	"regexp"
)

type webApp struct {
	savedGoalsPath  string
	webPath         string
	fileNameMatcher *regexp.Regexp
}

func (me *webApp) start() {
	if len(me.savedGoalsPath) == 0 {
		me.savedGoalsPath = "./saved-goals"
	}
	if me.fileNameMatcher == nil {
		me.fileNameMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d`)
	}
	http.HandleFunc(me.webPath+"/api/goals", me.wrap(me.getGoals))
}

func (me *webApp) wrap(function func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost:1234")
		function(response, request)
	}
}

func (me *webApp) getGoals(response http.ResponseWriter, request *http.Request) {
	var files = assertResultError(os.ReadDir("./saved-goals"))
	var headers = make([]*goalHeaderExtended, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			var headerFilePath = me.savedGoalsPath + "/" + file.Name() + "/_header.json"
			var header = readJsonFile(headerFilePath, &goalHeaderExtended{})
			me.extendHeader(header)
			headers = append(headers, header)
		}
	}
	response.Write(encodeJson(headers))
}

func (me *webApp) extendHeader(theGoalHeader *goalHeaderExtended) {
	var files = assertResultError(os.ReadDir(me.savedGoalsPath + "/" + theGoalHeader.Id))
	sortFilesByName(files)
	for i := len(files) - 1; i >= 0; i-- {
		if me.fileNameMatcher.MatchString(files[i].Name()) {
			var lastFileName = files[i].Name()
			theGoalHeader.LastPostDate = lastFileName[:len("2025-01-02")]
			break
		}
	}
}
