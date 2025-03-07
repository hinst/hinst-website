package main

import (
	"net/http"
	"os"
	"regexp"
	"time"
)

type webApp struct {
	savedGoalsPath        string
	webPath               string
	fileNameMatcher       *regexp.Regexp
	goalIdStringMatcher   *regexp.Regexp
	goalDateStringMatcher *regexp.Regexp
}

func (me *webApp) start() {
	me.savedGoalsPath = "./saved-goals"
	me.fileNameMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d`)
	me.goalIdStringMatcher = regexp.MustCompile(`^\d{1,10}$`)
	me.goalDateStringMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d$`)

	http.HandleFunc(me.webPath+"/api/goals", me.wrap(me.getGoals))
	http.HandleFunc(me.webPath+"/api/goalPosts", me.wrap(me.getGoalPosts))
	http.HandleFunc(me.webPath+"/api/goalPost", me.wrap(me.getGoalPost))
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
	var goalId = me.readValidGoalIdString(request.URL.Query().Get("id"))
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

func (me *webApp) getGoalPost(response http.ResponseWriter, request *http.Request) {
	var goalId = me.readValidGoalIdString(request.URL.Query().Get("goalId"))
	var postDateTime = me.readValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var postFileName = goalId + "/" + postDateTime.Format("2006-01-02_15-04-05") + ".json"
	var post = readJsonFile(me.savedGoalsPath+"/"+postFileName, &smartPost{})
	response.Write(encodeJson(post))
}

func (me *webApp) checkValidGoalIdString(goalId string) bool {
	return me.goalIdStringMatcher.MatchString(goalId)
}

func (me *webApp) readValidGoalIdString(goalId string) string {
	assertCondition(
		me.checkValidGoalIdString(goalId),
		webError{"Need valid goal id. Received: " + goalId, http.StatusBadRequest})
	return goalId
}

func (me *webApp) readValidPostDateTime(text string) time.Time {
	var postDateTime, postDateTimeError = parseSmartProgressDateTime(text)
	assertCondition(nil == postDateTimeError,
		webError{"Need valid postDateTime. Format: " + SMART_PROGRESS_DATE_TIME_FORMAT, http.StatusBadRequest})
	return postDateTime
}
