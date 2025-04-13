package main

import (
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"golang.org/x/text/language"
)

type webAppGoals struct {
	webAppGoalsBase
}

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
	var records = me.db.getGoals()
	response.Write(encodeJson(records))
}

func (me *webAppGoals) getGoal(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("id"))
	var goal = me.db.getGoal(goalId)
	setCacheAge(response, time.Minute)
	response.Write(encodeJson(goal))
}

func (me *webAppGoals) getGoalPosts(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("id"))
	var goalManagerMode = me.inputCheckGoalManagerMode(request)
	var posts = me.db.getGoalPosts(getIntFromString(goalId), goalManagerMode)
	response.Write(encodeJson(posts))
}

func (me *webAppGoals) getGoalPost(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var fileName = filepath.Join(goalId, postDateTime.Format(storedGoalFileTimeFormat)+".json")
	var filePath = filepath.Join(me.savedGoalsPath, fileName)
	var goalPost = readJsonFile(filePath, &smartPostExtended{})
	var goalManagerMode = me.inputCheckGoalManagerMode(request)

	var requestedLanguage = getWebLanguage(request)
	var translatedFilePath = translatorPresets.getTranslatedFilePath(filePath, requestedLanguage)
	if checkFileExists(translatedFilePath) {
		goalPost.Msg = readTextFile(translatedFilePath)
		goalPost.LanguageTag = requestedLanguage.String()
		goalPost.LanguageName = getLanguageName(requestedLanguage)
		if requestedLanguage != language.Russian {
			goalPost.IsAutoTranslated = true
		}
	} else {
		goalPost.LanguageNamePending = getLanguageName(requestedLanguage)
	}
	if goalManagerMode {
		var goalPostRow = me.db.getGoalPost(getIntFromString(goalId), postDateTime)
		if goalPostRow != nil {
			goalPost.IsPublic = goalPostRow.isPublic
		}
	}

	goalPost.Images = nil
	response.Write(encodeJson(goalPost))
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
	var isAdmin = me.inputCheckAdminPassword(request)
	if !isAdmin {
		panic(webError{"Need administrator access", http.StatusUnauthorized})
	}
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var isPublic = request.URL.Query().Get("isPublic") == "true"
	var row = goalPostRow{goalId: getIntFromString(goalId), dateTime: postDateTime, isPublic: isPublic}
	me.db.setGoalPostPublic(&row)
}
