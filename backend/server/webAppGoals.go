package server

import (
	"io"
	"net/http"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"github.com/hinst/hinst-website/server/rest_objects"
	"golang.org/x/text/language"
)

type webAppGoals struct {
	webAppGoalsBase
}

func (me *webAppGoals) init(db *database) []namedWebFunction {
	me.db = db
	return []namedWebFunction{
		{"/api/goals", me.getGoals},
		{"/api/goal/image", me.getGoalImage},
		{"/api/goal", me.getGoal},
		{"/api/goalPosts", me.getGoalPosts},
		{"/api/goalPost", me.getGoalPost},
		{"/api/goalPost/image", me.getGoalPostImage},
		{"/api/goalPost/setPublic", me.guardAdminFunction(me.setGoalPostPublic)},
		{"/api/goalPost/setSearchIndexingEnabled", me.guardAdminFunction(me.setGoalPostSearchIndexingEnabled)},
		{"/api/goalPost/setText", me.guardAdminFunction(me.setGoalPostText)},
		{"/api/goalPost/setTitle", me.guardAdminFunction(me.setGoalTitleText)},
		{"/api/goalPosts/search", me.searchGoalPosts},
	}
}

func (me *webAppGoals) getGoals(response http.ResponseWriter, request *http.Request) {
	var rows = me.db.getGoals()
	var records = []rest_objects.GoalObject{}
	for _, row := range rows {
		records = append(records, rest_objects.GoalObject{}.Read(row))
	}
	writeJsonResponse(response, records)
}

func (me *webAppGoals) getGoal(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("id"))
	var row = me.db.getGoal(goalId)
	if row == nil {
		panic(webError{"Goal not found", http.StatusNotFound})
	}
	writeJsonResponse(response, rest_objects.GoalObject{}.Read(*row))
}

func (me *webAppGoals) getGoalImage(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("id"))
	var imageData, imageContentType = me.db.getGoalImage(goalId)
	gophers.SetCacheAge(response, time.Hour)
	response.Header().Set(gophers.ContentTypeHeader, imageContentType)
	var _, _ = response.Write(imageData)
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
		var errorMessage = "Cannot find goalId=" + gophers.GetStringFromInt64(goalId) +
			" postDateTime=" + postDateTime.String()
		panic(webError{errorMessage, http.StatusNotFound})
	}
	if !goalPostRow.IsPublic && !goalManagerMode {
		panic(webError{"Need goal manager access level", http.StatusUnauthorized})
	}
	var goalPostObject rest_objects.GoalPostObject
	goalPostObject.GoalId = goalPostRow.GoalId
	goalPostObject.DateTime = goalPostRow.GetDateTime().UTC().Unix()
	goalPostObject.Text = goalPostRow.Text
	var requestedLanguage = getWebLanguage(request)
	goalPostObject.LanguageTag = requestedLanguage.String()
	goalPostObject.LanguageName = base.GetLanguageName(requestedLanguage)
	if requestedLanguage != base.SupportedLanguages[0] {
		var translatedText = goalPostRow.GetTranslatedText(requestedLanguage)
		if translatedText != "" {
			goalPostObject.IsAutoTranslated = true
			goalPostObject.Text = translatedText
		} else {
			goalPostObject.IsTranslationPending = true
		}
	}
	goalPostObject.IsPublic = goalPostRow.IsPublic
	if me.inputCheckGoalManagerMode(request) {
		goalPostObject.SearchIndexingEnabled = goalPostRow.SearchIndexingEnabled
	}
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
	gophers.SetCacheAge(response, time.Hour)
	response.Header().Set(gophers.ContentTypeHeader, image.ContentType)
	var _, _ = response.Write(image.File)
}

func (me *webAppGoals) setGoalPostPublic(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var isPublic = request.URL.Query().Get("isPublic") == "true"
	var row = db_objects.GoalPostRow{GoalId: goalId, DateTime: postDateTime.UTC().Unix(), IsPublic: isPublic}
	me.db.setGoalPostPublic(row)
}

func (me *webAppGoals) setGoalPostSearchIndexingEnabled(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var enabled = request.URL.Query().Get("enabled") == "true"
	var row = db_objects.GoalPostRow{GoalId: goalId, DateTime: postDateTime.UTC().Unix(), SearchIndexingEnabled: enabled}
	me.db.setGoalPostSearchIndexingEnabled(row)
}

func (me *webAppGoals) setGoalPostText(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var languageTagText = request.URL.Query().Get("languageTag")
	var languageTag = gophers.AssertResultError(language.Parse(languageTagText))
	var text = string(gophers.AssertResultError(io.ReadAll(request.Body)))
	me.db.setGoalPostText(goalId, postDateTime, languageTag, text)
}

func (me *webAppGoals) setGoalTitleText(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var languageTagText = request.URL.Query().Get("languageTag")
	var languageTag = gophers.AssertResultError(language.Parse(languageTagText))
	var text = string(gophers.AssertResultError(io.ReadAll(request.Body)))
	me.db.setGoalPostTitle(goalId, postDateTime, languageTag, text)
}

func (me *webAppGoals) searchGoalPosts(response http.ResponseWriter, request *http.Request) {
	const resultLimit = 100
	var queryText = request.URL.Query().Get("query")
	var requestedLanguage = getWebLanguage(request)
	var goalManagerMode = me.inputCheckGoalManagerMode(request)
	var rows = me.db.searchGoalPosts(queryText, requestedLanguage, goalManagerMode, resultLimit)
	var records []rest_objects.GoalPostHeader
	for _, row := range rows {
		var record rest_objects.GoalPostHeader
		record.GoalId = row.GoalId
		record.DateTime = row.GetDateTime().UTC().Unix()
		record.IsPublic = row.IsPublic
		if me.inputCheckGoalManagerMode(request) {
			record.SearchIndexingEnabled = row.SearchIndexingEnabled
		}
		record.Type = row.TypeString
		record.Title = row.GetTranslatedTitle(requestedLanguage)
		records = append(records, record)
	}
	writeJsonResponse(response, records)
}
