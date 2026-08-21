package server

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
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
	huma.Get(webApi, "/api/goalPost", me.getGoalPost)
	huma.Get(webApi, "/api/goalPost/image", me.getGoalPostImage)
	huma.Put(webApi, "/api/goalPost/setPublic", me.setGoalPostPublic)

	return []namedWebFunction{
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
	var posts = me.db.getGoalPosts(input.Id, webContext.isAdminMode(ctx), webContext.getLanguage(ctx))
	return rest_objects.NewSimpleResponse(posts), nil
}

func (me *webAppGoals) getGoalPost(ctx context.Context, input *struct {
	GoalId       int64 `query:"goalId"`
	PostDateTime int64 `query:"postDateTime"`
}) (*rest_objects.Response[*rest_objects.GoalPostObject], error) {
	var postDateTime = time.Unix(input.PostDateTime, 0)
	var goalPostRow = me.db.getGoalPost(input.GoalId, postDateTime)
	var isNotFound = goalPostRow == nil ||
		!goalPostRow.IsPublic && !webContext.isAdminMode(ctx)
	if isNotFound {
		var errorMessage = "Cannot find goalId=" + gophers.GetStringFromInt64(input.GoalId) +
			" postDateTime=" + postDateTime.String()
		panic(webError{errorMessage, http.StatusNotFound})
	}
	var requestedLanguage = webContext.getLanguage(ctx)
	var goalPostObject rest_objects.GoalPostObject
	goalPostObject.GoalId = goalPostRow.GoalId
	goalPostObject.DateTime = goalPostRow.GetDateTime().UTC().Unix()
	goalPostObject.Text = goalPostRow.Text
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
	if webContext.isAdminMode(ctx) {
		goalPostObject.SearchIndexingEnabled = goalPostRow.SearchIndexingEnabled
	}
	goalPostObject.ImageCount = me.db.getGoalPostImageCount(input.GoalId, postDateTime)
	return rest_objects.NewSimpleResponse(&goalPostObject), nil
}

func (me *webAppGoals) getGoalPostImage(ctx context.Context, input *struct {
	GoalId       int64 `query:"goalId"`
	PostDateTime int64 `query:"postDateTime"`
	Index        int   `query:"index"`
}) (*rest_objects.Response[[]byte], error) {
	var postDateTime = time.Unix(input.PostDateTime, 0)
	var image = me.db.getGoalPostImage(input.GoalId, postDateTime, input.Index)
	if image == nil {
		panic(webError{"Image not found", http.StatusNotFound})
	}
	return &rest_objects.Response[[]byte]{
		Body:         image.File,
		ContentType:  image.ContentType,
		CacheControl: "max-age=" + strconv.Itoa(int(time.Hour.Seconds())),
	}, nil
}

func (me *webAppGoals) setGoalPostPublic(ctx context.Context, input *struct {
	GoalId       int64 `query:"goalId" required:"true"`
	PostDateTime int64 `query:"postDateTime" required:"true"`
	IsPublic     bool  `query:"isPublic" required:"true"`
}) (*struct{}, error) {
	if !webContext.isAdminMode(ctx) {
		panic(webError{"Need admin password", http.StatusForbidden})
	}
	var postDateTime = time.Unix(input.PostDateTime, 0)
	var row = db_objects.GoalPostRow{GoalId: input.GoalId, DateTime: postDateTime.UTC().Unix(), IsPublic: input.IsPublic}
	me.db.setGoalPostPublic(row)
	return &struct{}{}, nil
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
