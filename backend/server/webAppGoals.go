package server

import (
	"context"
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

func (me *webAppGoals) init(webApi huma.API, db *database) {
	me.db = db
	huma.Get(webApi, "/api/goals", me.getGoals)
	huma.Get(webApi, "/api/goal", me.getGoal)
	huma.Get(webApi, "/api/goal/image", me.getGoalImage)
	huma.Get(webApi, "/api/goalPosts", me.getGoalPosts)
	huma.Get(webApi, "/api/goalPost", me.getGoalPost)
	huma.Get(webApi, "/api/goalPost/image", me.getGoalPostImage)
	huma.Put(webApi, "/api/goalPost/setPublic", me.setGoalPostPublic)
	huma.Put(webApi, "/api/goalPost/setSearchIndexingEnabled", me.setGoalPostSearchIndexingEnabled)
	huma.Post(webApi, "/api/goalPost/setText", me.setGoalPostText)
	huma.Post(webApi, "/api/goalPost/setTitle", me.setGoalTitleText)
	huma.Get(webApi, "/api/goalPosts/search", me.searchGoalPosts)
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
	Id int64 `query:"id" required:"true"`
}) (*rest_objects.Response[*rest_objects.GoalObject], error) {
	var row = me.db.getGoal(input.Id)
	if row == nil {
		panic(webError{"Goal not found", http.StatusNotFound})
	}
	var goalObject = rest_objects.GoalObject{}.Read(*row)
	return rest_objects.NewSimpleResponse(&goalObject), nil
}

func (me *webAppGoals) getGoalImage(ctx context.Context, input *struct {
	Id int64 `query:"id" required:"true"`
}) (*rest_objects.Response[[]byte], error) {
	var imageData, imageContentType = me.db.getGoalImage(input.Id)
	return &rest_objects.Response[[]byte]{
		Body:         imageData,
		ContentType:  imageContentType,
		CacheControl: "max-age=" + strconv.Itoa(int(time.Hour.Seconds())),
	}, nil
}

func (me *webAppGoals) getGoalPosts(ctx context.Context, input *struct {
	Id int64 `query:"id" required:"true"`
}) (*rest_objects.Response[[]rest_objects.GoalPostHeader], error) {
	var posts = me.db.getGoalPosts(input.Id, webContext.isAdminMode(ctx), webContext.getLanguage(ctx))
	return rest_objects.NewSimpleResponse(posts), nil
}

func (me *webAppGoals) getGoalPost(ctx context.Context, input *struct {
	GoalId       int64 `query:"goalId" required:"true"`
	PostDateTime int64 `query:"postDateTime" required:"true"`
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
	GoalId       int64 `query:"goalId" required:"true"`
	PostDateTime int64 `query:"postDateTime" required:"true"`
	Index        int   `query:"index" required:"true"`
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
	webContext.assertAdminMode(ctx)
	var postDateTime = time.Unix(input.PostDateTime, 0)
	var row = db_objects.GoalPostRow{GoalId: input.GoalId, DateTime: postDateTime.UTC().Unix(), IsPublic: input.IsPublic}
	me.db.setGoalPostPublic(row)
	return &struct{}{}, nil
}

func (me *webAppGoals) setGoalPostSearchIndexingEnabled(ctx context.Context, input *struct {
	GoalId       int64 `query:"goalId" required:"true"`
	PostDateTime int64 `query:"postDateTime" required:"true"`
	Enabled      bool  `query:"enabled" required:"true"`
}) (*struct{}, error) {
	webContext.assertAdminMode(ctx)
	var postDateTime = time.Unix(input.PostDateTime, 0)
	var row = db_objects.GoalPostRow{GoalId: input.GoalId, DateTime: postDateTime.UTC().Unix(), SearchIndexingEnabled: input.Enabled}
	me.db.setGoalPostSearchIndexingEnabled(row)
	return &struct{}{}, nil
}

func (me *webAppGoals) setGoalPostText(ctx context.Context, input *struct {
	GoalId       int64  `query:"goalId" required:"true"`
	PostDateTime int64  `query:"postDateTime" required:"true"`
	LanguageTag  string `query:"languageTag" required:"true"`
	RawBody      []byte `contentType:"text/plain;charset=UTF-8"` // The name of this parameter needs to be exactly RawBody
}) (*struct{}, error) {
	webContext.assertAdminMode(ctx)
	var languageTag, parseError = language.Parse(input.LanguageTag)
	gophers.AssertCondition(parseError == nil, func() webError {
		return webError{"Need valid language tag. Received: " + input.LanguageTag, http.StatusBadRequest}
	})
	var text = string(input.RawBody)
	me.db.setGoalPostText(input.GoalId, time.Unix(input.PostDateTime, 0), languageTag, text)
	return &struct{}{}, nil
}

func (me *webAppGoals) setGoalTitleText(ctx context.Context, input *struct {
	GoalId       int64  `query:"goalId" required:"true"`
	PostDateTime int64  `query:"postDateTime" required:"true"`
	LanguageTag  string `query:"languageTag" required:"true"`
	RawBody      []byte `contentType:"text/plain;charset=UTF-8"`
}) (*struct{}, error) {
	webContext.assertAdminMode(ctx)
	var languageTag, parseError = language.Parse(input.LanguageTag)
	gophers.AssertCondition(parseError == nil, func() webError {
		return webError{"Need valid language tag. Received: " + input.LanguageTag, http.StatusBadRequest}
	})
	var text = string(input.RawBody)
	me.db.setGoalPostTitle(input.GoalId, time.Unix(input.PostDateTime, 0), languageTag, text)
	return &struct{}{}, nil
}

func (me *webAppGoals) searchGoalPosts(ctx context.Context, input *struct {
	Query string `query:"query"`
}) (*rest_objects.Response[[]rest_objects.GoalPostHeader], error) {
	const resultLimit = 100
	var rows = me.db.searchGoalPosts(input.Query, webContext.getLanguage(ctx),
		webContext.isAdminMode(ctx), resultLimit)
	var records []rest_objects.GoalPostHeader
	for _, row := range rows {
		var record rest_objects.GoalPostHeader
		record.Read(row, webContext.getLanguage(ctx))
		records = append(records, record)
	}
	return rest_objects.NewSimpleResponse(records), nil
}
