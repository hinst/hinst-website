package server

import (
	"context"
	"net/http"

	"github.com/hinst/go-gophers"
	"golang.org/x/text/language"
)

type tWebContext struct {
	keyLanguage string
	keyIsAdmin  string
}

var webContext = tWebContext{
	keyLanguage: "webContextKeyLanguage",
	keyIsAdmin:  "webContextKeyIsAdmin",
}

func (me tWebContext) getLanguage(ctx context.Context) language.Tag {
	return ctx.Value(me.keyLanguage).(language.Tag)
}

func (me tWebContext) isAdminMode(ctx context.Context) bool {
	return ctx.Value(me.keyIsAdmin).(bool)
}

func (me tWebContext) assertAdminMode(ctx context.Context) {
	gophers.AssertCondition(me.isAdminMode(ctx), func() webError {
		return webError{Message: "Need admin mode", Status: http.StatusUnauthorized}
	})
}
