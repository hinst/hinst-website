package main

import (
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"golang.org/x/text/language"
)

type webAppGoals struct {
	webAppGoalsBase
}

func (me *webAppGoals) init(db *database) []namedWebFunction {
	me.db = db
	return []namedWebFunction{
		{"/api/goals", me.getGoals},
		{"/api/goal", me.getGoal},
		{"/api/goalPosts", me.getGoalPosts},
		{"/api/goalPost", me.getGoalPost},
		{"/api/goalPost/images", me.getGoalPostImages},
		{"/api/goalPost/setPublic", me.setGoalPostPublic},
		{"/api/goalPost/setText", me.setGoalPostText},
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
	var goalId = me.inputValidGoalId(request.URL.Query().Get("id"))
	var goalManagerMode = me.inputCheckGoalManagerMode(request)
	var posts = me.db.getGoalPosts(goalId, goalManagerMode)
	response.Write(encodeJson(posts))
}

func (me *webAppGoals) getGoalPost(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var goalManagerMode = me.inputCheckGoalManagerMode(request)

	var goalPostRow = me.db.getGoalPost(goalId, postDateTime)
	if goalPostRow == nil {
		var errorMessage = "Cannot find goalId=" + getStringFromInt64(goalId) +
			" postDateTime=" + postDateTime.String()
		panic(webError{errorMessage, http.StatusNotFound})
	}
	if !goalPostRow.isPublic && !goalManagerMode {
		panic(webError{"Need goal manager access level", http.StatusUnauthorized})
	}
	var goalPostObject goalPostObject
	goalPostObject.GoalId = goalPostRow.goalId
	goalPostObject.DateTime = goalPostRow.dateTime.UTC().Unix()
	goalPostObject.Text = goalPostRow.text
	var requestedLanguage = getWebLanguage(request)
	goalPostObject.LanguageTag = requestedLanguage.String()
	goalPostObject.LanguageName = getLanguageName(requestedLanguage)
	if requestedLanguage != supportedLanguages[0] {
		var translatedText = goalPostRow.getTranslatedText(requestedLanguage)
		if translatedText != "" {
			goalPostObject.IsAutoTranslated = true
			goalPostObject.Text = translatedText
		} else {
			goalPostObject.IsTranslationPending = true
		}
	}
	goalPostObject.IsPublic = goalPostRow.isPublic
	response.Write(encodeJson(goalPostObject))
}

func (me *webAppGoals) getGoalPostImages(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var images = me.db.getGoalPostImages(goalId, postDateTime)
	var imageStrings []string
	for _, image := range images {
		var imageString = "data:" + image.contentType + ";base64," +
			base64.StdEncoding.EncodeToString(image.file)
		imageStrings = append(imageStrings, imageString)
	}
	setCacheAge(response, time.Minute)
	response.Write(encodeJson(imageStrings))
}

func (me *webAppGoals) setGoalPostPublic(response http.ResponseWriter, request *http.Request) {
	var isAdmin = me.inputCheckAdminPassword(request)
	if !isAdmin {
		panic(webError{"Need administrator access", http.StatusUnauthorized})
	}
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var isPublic = request.URL.Query().Get("isPublic") == "true"
	var row = goalPostRow{goalId: goalId, dateTime: postDateTime, isPublic: isPublic}
	me.db.setGoalPostPublic(row)
}

func (me *webAppGoals) setGoalPostText(response http.ResponseWriter, request *http.Request) {
	var isAdmin = me.inputCheckAdminPassword(request)
	if !isAdmin {
		panic(webError{"Need administrator access", http.StatusUnauthorized})
	}
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var languageTagText = request.URL.Query().Get("languageTag")
	var text = string(assertResultError(io.ReadAll(request.Body)))
	var languageTag = assertResultError(language.Parse(languageTagText))
	me.db.setGoalPostText(goalId, postDateTime, languageTag, text)
}
