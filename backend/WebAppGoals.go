package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type webAppGoals struct {
	savedGoalsPath        string
	goalIdStringMatcher   *regexp.Regexp
	goalDateStringMatcher *regexp.Regexp
}

const savedGoalHeaderFileName = "_header.json"

func (me *webAppGoals) init() []namedWebFunction {
	me.goalIdStringMatcher = regexp.MustCompile(`^\d{1,10}$`)
	me.goalDateStringMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d$`)
	return []namedWebFunction{
		{"/api/goals", me.getGoals},
		{"/api/goal", me.getGoal},
		{"/api/goalPosts", me.getGoalPosts},
		{"/api/goalPost", me.getGoalPost},
	}
}

func (me *webAppGoals) getGoals(response http.ResponseWriter, request *http.Request) {
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

func (me *webAppGoals) getGoal(response http.ResponseWriter, request *http.Request) {
	var goalId = me.readValidGoalIdString(request.URL.Query().Get("id"))
	var headerFilePath = filepath.Join(me.savedGoalsPath, goalId, savedGoalHeaderFileName)
	var goal = readJsonFile(headerFilePath, &goalHeader{})
	response.Write(encodeJson(goal))
}

func (me *webAppGoals) extendHeader(goalHeader *goalHeaderExtended) {
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

func (me *webAppGoals) getGoalPosts(response http.ResponseWriter, request *http.Request) {
	var goalId = me.readValidGoalIdString(request.URL.Query().Get("id"))
	var goalDirectoryPath = filepath.Join(me.savedGoalsPath, goalId)
	var files = assertResultError(os.ReadDir(goalDirectoryPath))
	sortFilesByName(files)
	var posts = make([]*smartPostHeader, 0, len(files))
	for _, file := range files {
		if GoalFileNameMatcher.MatchString(file.Name()) {
			var filePath = filepath.Join(me.savedGoalsPath, goalId, file.Name())
			var post = readJsonFile(filePath, &smartPostHeader{})
			posts = append(posts, post)
		}
	}
	response.Write(encodeJson(posts))
}

func (me *webAppGoals) getGoalPost(response http.ResponseWriter, request *http.Request) {
	var goalId = me.readValidGoalIdString(request.URL.Query().Get("goalId"))
	var postDateTime = me.readValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var fileName = filepath.Join(goalId, postDateTime.Format(storedGoalFileTimeFormat)+".json")
	var filePath = filepath.Join(me.savedGoalsPath, fileName)
	var post = readJsonFile(filePath, &smartPost{})

	var requestedLanguage = getWebLanguage(request)
	log.Printf("Requested language: %v", requestedLanguage)
	var translatedFilePath = translatorPresets.getTranslatedFilePath(filePath, requestedLanguage)
	if checkFileExists(translatedFilePath) {
		post.Msg = readTextFile(translatedFilePath)
	}

	response.Write(encodeJson(post))
}

func (me *webAppGoals) checkValidGoalIdString(goalId string) bool {
	return me.goalIdStringMatcher.MatchString(goalId)
}

func (me *webAppGoals) readValidGoalIdString(goalId string) string {
	var createWebError = func() webError {
		return webError{"Need goal id. Received: " + goalId, http.StatusBadRequest}
	}
	assertCondition(me.checkValidGoalIdString(goalId), createWebError)
	return goalId
}

func (me *webAppGoals) readValidPostDateTime(text string) time.Time {
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
