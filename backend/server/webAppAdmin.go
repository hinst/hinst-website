package server

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"github.com/hinst/hinst-website/server/rest_objects"
)

type webAppAdmin struct {
	webAppBase
	db *database
}

func (me *webAppAdmin) init(webApi huma.API, db *database) {
	me.db = db
	huma.Get(webApi, "/api/urlPings", me.getUrlPings)
}

func (me *webAppAdmin) getUrlPings(ctx context.Context, input *struct{}) (*rest_objects.Response[[]rest_objects.GoalPostSearchIndexingHeader], error) {
	if !webContext.isAdminMode(ctx) {
		panic(webError{"Need admin mode", http.StatusForbidden})
	}
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
