package main

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type webApp struct {
	savedGoalsPath        string
	webPath               string
	goalIdStringMatcher   *regexp.Regexp
	goalDateStringMatcher *regexp.Regexp
}

const savedGoalHeaderFileName = "_header.json"

func (me *webApp) init() {
	if me.webPath == "" {
		me.webPath = "/hinst-website"
	}
	me.goalIdStringMatcher = regexp.MustCompile(`^\d{1,10}$`)
	me.goalDateStringMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d$`)

	http.HandleFunc(me.webPath+"/api/goals", me.wrap(me.getGoals))
	http.HandleFunc(me.webPath+"/api/goal", me.wrap(me.getGoal))
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
	var files = assertResultError(os.ReadDir(me.savedGoalsPath))
	var headers = make([]*goalHeaderExtended, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			var headerFilePath = filepath.Join(me.savedGoalsPath, file.Name(), savedGoalHeaderFileName)
			var header = readJsonFile(headerFilePath, &goalHeaderExtended{})
			me.extendHeader(header)
			headers = append(headers, header)
		}
	}
	response.Write(encodeJson(headers))
}

func (me *webApp) getGoal(response http.ResponseWriter, request *http.Request) {
	var goalId = me.readValidGoalIdString(request.URL.Query().Get("id"))
	var headerFilePath = filepath.Join(me.savedGoalsPath, goalId, savedGoalHeaderFileName)
	var goal = readJsonFile(headerFilePath, &goalHeader{})
	response.Write(encodeJson(goal))
}

func (me *webApp) extendHeader(goalHeader *goalHeaderExtended) {
	var goalDirectoryPath = filepath.Join(me.savedGoalsPath, goalHeader.Id)
	var files = assertResultError(os.ReadDir(goalDirectoryPath))
	sortFilesByName(files)
	for i := len(files) - 1; i >= 0; i-- {
		if GoalFileNameMatcher.MatchString(files[i].Name()) {
			var lastFileName = files[i].Name()
			goalHeader.LastPostDate = lastFileName[:len("2025-01-02")]
			break
		}
	}
}

func (me *webApp) getGoalPosts(response http.ResponseWriter, request *http.Request) {
	var goalId = me.readValidGoalIdString(request.URL.Query().Get("id"))
	var goalDirectoryPath = filepath.Join(me.savedGoalsPath, goalId)
	var files = assertResultError(os.ReadDir(goalDirectoryPath))
	sortFilesByName(files)
	var posts = make([]*smartPostHeader, 0, len(files))
	for _, file := range files {
		if GoalFileNameMatcher.MatchString(file.Name()) {
			var post = readJsonFile(me.savedGoalsPath+"/"+goalId+"/"+file.Name(), &smartPostHeader{})
			posts = append(posts, post)
		}
	}
	response.Write(encodeJson(posts))
}

func (me *webApp) getGoalPost(response http.ResponseWriter, request *http.Request) {
	var goalId = me.readValidGoalIdString(request.URL.Query().Get("goalId"))
	var postDateTime = me.readValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var postFileName = filepath.Join(goalId, postDateTime.Format(storedGoalFileTimeFormat)+".json")
	var post = readJsonFile(me.savedGoalsPath+"/"+postFileName, &smartPost{})
	response.Write(encodeJson(post))
}

func (me *webApp) checkValidGoalIdString(goalId string) bool {
	return me.goalIdStringMatcher.MatchString(goalId)
}

func (me *webApp) readValidGoalIdString(goalId string) string {
	var createWebError = func() webError {
		return webError{"Need goal id. Received: " + goalId, http.StatusBadRequest}
	}
	assertCondition(me.checkValidGoalIdString(goalId), createWebError)
	return goalId
}

func (me *webApp) readValidPostDateTime(text string) time.Time {
	var postDateTime, postDateTimeError = parseSmartProgressDateTime(text)
	var createWebError = func() webError {
		return webError{
			"Need valid postDateTime. Format: " + smartProgressTimeFormat,
			http.StatusBadRequest,
		}
	}
	assertCondition(nil == postDateTimeError, createWebError)
	return postDateTime
}
