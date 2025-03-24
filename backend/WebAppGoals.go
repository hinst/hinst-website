package main

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"golang.org/x/text/language"
)

type webAppGoals struct {
	webAppGoalsBase
}

const cookieKeyAdminPassword = "adminPassword"

func (me *webAppGoals) init(db *database) []namedWebFunction {
	me.db = db
	me.goalDateStringMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d$`)
	return []namedWebFunction{
		{"/api/goals", me.getGoals},
		{"/api/goal", me.getGoal},
		{"/api/goalPosts", me.getGoalPosts},
		{"/api/goalPost", me.getGoalPost},
		{"/api/goalPost/images", me.getGoalPostImages},
		{"/api/goalPost/setPublic", me.setGoalPostPublic},
	}
}

func (me *webAppGoals) getGoals(response http.ResponseWriter, request *http.Request) {
	var files = assertResultError(os.ReadDir(me.savedGoalsPath))
	var headers = make([]*goalHeader, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			var headerFilePath = filepath.Join(me.savedGoalsPath, file.Name(), savedGoalHeaderFileName)
			var header = readJsonFile(headerFilePath, &goalHeader{})
			headers = append(headers, header)
		}
	}
	response.Write(encodeJson(headers))
}

func (me *webAppGoals) getGoal(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("id"))
	var headerFilePath = filepath.Join(me.savedGoalsPath, goalId, savedGoalHeaderFileName)
	var goal = readJsonFile(headerFilePath, &goalHeader{})
	setCacheAge(response, time.Minute)
	response.Write(encodeJson(goal))
}

func (me *webAppGoals) getGoalPosts(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("id"))
	var goalManagerMode = me.inputCheckGoalManagerMode(request)
	var filePaths = me.getGoalPostFiles(goalId, goalManagerMode)
	var posts = readJsonFiles[smartPostHeader](filePaths, runtime.NumCPU())
	setCacheAge(response, time.Minute)
	response.Write(encodeJson(posts))
}

func (me *webAppGoals) getGoalPost(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var fileName = filepath.Join(goalId, postDateTime.Format(storedGoalFileTimeFormat)+".json")
	var filePath = filepath.Join(me.savedGoalsPath, fileName)
	var post = readJsonFile(filePath, &smartPostExtended{})
	var goalManagerMode = me.inputCheckGoalManagerMode(request)

	var requestedLanguage = getWebLanguage(request)
	var translatedFilePath = translatorPresets.getTranslatedFilePath(filePath, requestedLanguage)
	if checkFileExists(translatedFilePath) {
		post.Msg = readTextFile(translatedFilePath)
		post.LanguageTag = requestedLanguage.String()
		post.LanguageName = getLanguageName(requestedLanguage)
		if requestedLanguage != language.Russian {
			post.IsAutoTranslated = true
		}
	} else {
		post.LanguageNamePending = getLanguageName(requestedLanguage)
	}
	if goalManagerMode {
		post.IsPublic = me.db.getAvailablePosts(getIntFromString(goalId))[post.Date]
	}

	post.Images = nil
	response.Write(encodeJson(post))
}

func (me *webAppGoals) getGoalPostImages(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var fileName = filepath.Join(goalId, postDateTime.Format(storedGoalFileTimeFormat)+".json")
	var filePath = filepath.Join(me.savedGoalsPath, fileName)
	var post = readJsonFile(filePath, &smartPostExtended{})
	setCacheAge(response, time.Minute)
	response.Write(encodeJson(post.Images))
}

func (me *webAppGoals) setGoalPostPublic(response http.ResponseWriter, request *http.Request) {
}
