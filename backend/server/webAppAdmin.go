package server

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"github.com/hinst/hinst-website/server/rest_objects"
)

type webAppAdmin struct {
	db *database
}

func (me *webAppAdmin) init(webApi huma.API, db *database) {
	me.db = db
	huma.Get(webApi, "/api/urlPings", me.getUrlPings)
	huma.Get(webApi, "/api/isAdminModeEnabled", me.getIsAdminModeEnabled)
	huma.Put(webApi, "/api/goalPosts/googlePingedAt", me.setGoalPostGooglePingedAt)
}

func (me *webAppAdmin) getIsAdminModeEnabled(ctx context.Context, input *struct{}) (*rest_objects.Response[bool], error) {
	return rest_objects.NewSimpleResponse(webContext.isAdminMode(ctx)), nil
}

func (me *webAppAdmin) getUrlPings(ctx context.Context, input *struct{}) (*rest_objects.Response[[]rest_objects.GoalPostSearchIndexingHeader], error) {
	webContext.assertAdminMode(ctx)
	var objects = []rest_objects.GoalPostSearchIndexingHeader{}
	var webLanguage = base.SupportedLanguages[0] // Only blog posts in main language currently participate in search indexing
	me.db.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		if !row.SearchIndexingEnabled {
			return true
		}
		var header rest_objects.GoalPostSearchIndexingHeader
		header.Read(row, webLanguage)
		header.PublicUrl = webStaticGoals{}.getPublicUrl(row, webLanguage)
		objects = append(objects, header)
		return true
	}, (db_objects.GoalPostRow{}).GetAllFieldSelector(), -1)
	return rest_objects.NewSimpleResponse(objects), nil
}

func (me *webAppAdmin) setGoalPostGooglePingedAt(ctx context.Context, input *struct {
	GoalId       int64 `query:"goalId" required:"true"`
	PostDateTime int64 `query:"postDateTime" required:"true"`
}) (*rest_objects.Response[bool], error) {
	webContext.assertAdminMode(ctx)
	var postDateTime = time.Unix(input.PostDateTime, 0)
	var affectedRowCount = me.db.setGoalPostGooglePingedAt(input.GoalId, postDateTime, time.Now())
	return rest_objects.NewSimpleResponse(affectedRowCount > 0), nil
}
