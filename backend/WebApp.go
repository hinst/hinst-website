package main

import (
	"net/http"
	"os"
	"regexp"
)

type webApp struct {
	savedGoalsPath      string
	webPath             string
	fileNameMatcher     *regexp.Regexp
	goalIdStringMatcher *regexp.Regexp
}

func (me *webApp) start() {
	me.savedGoalsPath = "./saved-goals"
	me.fileNameMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d`)
	me.goalIdStringMatcher = regexp.MustCompile(`^\d+$`)
	http.HandleFunc(me.webPath+"/api/goals", me.wrap(me.getGoals))
	http.HandleFunc(me.webPath+"/api/goalPosts", me.wrap(me.getGoalPosts))
}

func (me *webApp) wrap(function func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost:1234")
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

func (me *webApp) getGoalPosts(response http.ResponseWriter, request *http.Request) {
	var goalId = request.URL.Query().Get("goalId")
	assertCondition(
		me.checkValidGoalIdString(goalId),
		webError{"Need valid goal id", http.StatusBadRequest})
	var files = assertResultError(os.ReadDir(me.savedGoalsPath + "/" + goalId))
	sortFilesByName(files)
	var posts = make([]*smartPostHeader, 0, len(files))
	for _, file := range files {
		if me.fileNameMatcher.MatchString(file.Name()) {
			var post = readJsonFile(me.savedGoalsPath+"/"+goalId+"/"+file.Name(), &smartPostHeader{})
			posts = append(posts, post)
		}
	}
	response.Write(encodeJson(posts))
}

func (me *webApp) checkValidGoalIdString(goalId string) bool {
	return me.goalIdStringMatcher.MatchString(goalId)
}
