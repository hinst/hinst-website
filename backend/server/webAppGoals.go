package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"github.com/hinst/hinst-website/server/rest_objects"
	"golang.org/x/text/language"
)

type webAppGoals struct {
	webAppGoalsBase
}

func (me *webAppGoals) init(webApi huma.API, db *database) []namedWebFunction {
	me.db = db
	huma.Get(webApi, "/api/goals", me.getGoals)
	huma.Get(webApi, "/api/goal", me.getGoal)
	huma.Get(webApi, "/api/goal/image", me.getGoalImage)
	huma.Get(webApi, "/api/goalPosts", me.getGoalPosts)

	return []namedWebFunction{
		{"/api/goalPost", me.getGoalPost},
		{"/api/goalPost/image", me.getGoalPostImage},
		{"/api/goalPost/setPublic", me.guardAdminFunction(me.setGoalPostPublic)},
		{"/api/goalPost/setSearchIndexingEnabled", me.guardAdminFunction(me.setGoalPostSearchIndexingEnabled)},
		{"/api/goalPost/setText", me.guardAdminFunction(me.setGoalPostText)},
		{"/api/goalPost/setTitle", me.guardAdminFunction(me.setGoalTitleText)},
		{"/api/goalPosts/search", me.searchGoalPosts},
	}
}

func (me *webAppGoals) getGoals(ctx context.Context, input *struct{}) (*rest_objects.Response[[]rest_objects.GoalObject], error) {
	var rows = me.db.getGoals()
	var records = []rest_objects.GoalObject{}
	for _, row := range rows {
		records = append(records, rest_objects.GoalObject{}.Read(row))
	}
	return rest_objects.NewSimpleResponse(records), nil
}

func (me *webAppGoals) getGoal(ctx context.Context, input *struct {
	Id int64 `query:"id"`
}) (*rest_objects.Response[*rest_objects.GoalObject], error) {
	var row = me.db.getGoal(input.Id)
	if row == nil {
		panic(webError{"Goal not found", http.StatusNotFound})
	}
	var goalObject = rest_objects.GoalObject{}.Read(*row)
	return rest_objects.NewSimpleResponse(&goalObject), nil
}

func (me *webAppGoals) getGoalImage(ctx context.Context, input *struct {
	Id int64 `query:"id"`
}) (*rest_objects.Response[[]byte], error) {
	log.Println(input.Id)
	var imageData, imageContentType = me.db.getGoalImage(input.Id)
	return &rest_objects.Response[[]byte]{
		Body:         imageData,
		ContentType:  imageContentType,
		CacheControl: "max-age=" + strconv.Itoa(int(time.Hour.Seconds())),
	}, nil
}

func (me *webAppGoals) getGoalPosts(ctx context.Context, input *struct {
	Id int64 `query:"id"`
}) (*rest_objects.Response[[]rest_objects.GoalPostHeader], error) {
	var request, _ = humago.Unwrap(ctx)
	var goalManagerMode = me.inputCheckGoalManagerMode(request)
	var requestedLanguage = ctx.Value(webContextKeyLanguage).(language.Tag)
	var posts = me.db.getGoalPosts(input.Id, goalManagerMode, requestedLanguage)
	return rest_objects.NewSimpleResponse(posts), nil
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
		record.Read(row, requestedLanguage)
		records = append(records, record)
	}
	writeJsonResponse(response, records)
}
