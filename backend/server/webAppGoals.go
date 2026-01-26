package server

import (
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
		{"/api/goalPost/image", me.getGoalPostImage},
		{"/api/goalPost/setPublic", me.guardAdminFunction(me.setGoalPostPublic)},
		{"/api/goalPost/setText", me.guardAdminFunction(me.setGoalPostText)},
		{"/api/goalPost/setTitle", me.guardAdminFunction(me.setGoalTitleText)},
		{"/api/goalPost/search", me.searchGoalPosts},
	}
}

func (me *webAppGoals) getGoals(response http.ResponseWriter, request *http.Request) {
	var records = me.db.getGoals()
	writeJsonResponse(response, records)
}

func (me *webAppGoals) getGoal(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("id"))
	var goal = me.db.getGoal(goalId)
	setCacheAge(response, time.Minute)
	writeJsonResponse(response, goal)
}

func (me *webAppGoals) getGoalPosts(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("id"))
	var goalManagerMode = me.inputCheckGoalManagerMode(request)
	var requestedLanguage = getWebLanguage(request)
	var posts = me.db.getGoalPosts(goalId, goalManagerMode, requestedLanguage)
	writeJsonResponse(response, posts)
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
	goalPostObject.ImageCount = me.db.getGoalPostImageCount(goalId, postDateTime)
	writeJsonResponse(response, goalPostObject)
}

func (me *webAppGoals) getGoalPostImage(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var index = inputValidWebInteger(request.URL.Query().Get("index"))
	var image = me.db.getGoalPostImage(goalId, postDateTime, index)
	if image == nil {
		panic(webError{"Image not found", http.StatusNotFound})
	}
	setCacheAge(response, time.Hour)
	response.Header().Set("Content-Type", image.contentType)
	var _, _ = response.Write(image.file)
}

func (me *webAppGoals) setGoalPostPublic(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var isPublic = request.URL.Query().Get("isPublic") == "true"
	var row = goalPostRow{goalId: goalId, dateTime: postDateTime, isPublic: isPublic}
	me.db.setGoalPostPublic(row)
}

func (me *webAppGoals) setGoalPostText(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var languageTagText = request.URL.Query().Get("languageTag")
	var languageTag = assertResultError(language.Parse(languageTagText))
	var text = string(assertResultError(io.ReadAll(request.Body)))
	me.db.setGoalPostText(goalId, postDateTime, languageTag, text)
}

func (me *webAppGoals) setGoalTitleText(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var languageTagText = request.URL.Query().Get("languageTag")
	var languageTag = assertResultError(language.Parse(languageTagText))
	var text = string(assertResultError(io.ReadAll(request.Body)))
	me.db.setGoalPostTitle(goalId, postDateTime, languageTag, text)
}

func (me *webAppGoals) searchGoalPosts(response http.ResponseWriter, request *http.Request) {
	const resultLimit = 100
	var queryText = request.URL.Query().Get("query")
	var requestedLanguage = getWebLanguage(request)
	use(requestedLanguage)
	var goalManagerMode = me.inputCheckGoalManagerMode(request)
	var rows = me.db.searchGoalPosts(queryText, goalManagerMode, resultLimit)
	var records []goalPostRecord
	for _, row := range rows {
		var record goalPostRecord
		record.GoalId = row.goalId
		record.DateTime = row.dateTime.UTC().Unix()
		record.Type = row.typeString
		var title = row.getTranslatedTitle(requestedLanguage)
		record.Title = &title
		records = append(records, record)
	}
	writeJsonResponse(response, records)
}
