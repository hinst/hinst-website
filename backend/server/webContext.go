package server

import (
	"context"

	"golang.org/x/text/language"
)

const webContextKeyLanguage = "webContextKeyLanguage"

func getWebLanguageFromContext(ctx context.Context) language.Tag {
	return ctx.Value(webContextKeyLanguage).(language.Tag)
}

const webContextKeyIsAdmin = "webContextKeyIsAdmin"

func getWebAdminModeFromContext(ctx context.Context) bool {
	return ctx.Value(webContextKeyIsAdmin).(bool)
}
