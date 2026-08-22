package server

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/hinst/go-gophers"
	"golang.org/x/text/language"
)

type webApp struct {
	db          *database
	webPath     string
	webApi      huma.API
	webApiGroup huma.API
}

func (webApp) getDefaultWebPath() string {
	return "/hinst-website"
}

func (me *webApp) init(db *database) {
	me.db = db
	if me.webPath == "" {
		me.webPath = me.getDefaultWebPath()
	}
	if me.webPath == "/" {
		me.webPath = ""
	}
	var humaConfig = huma.DefaultConfig("Hinst-website API", "0.1.0")
	humaConfig.DocsPath = me.webPath + "/docs"
	me.webApi = humago.New(webRouter(), humaConfig)
	if me.webPath != "" {
		me.webApiGroup = huma.NewGroup(me.webApi, me.webPath)
	}
	me.webApi.UseMiddleware(me.catchPanic)
	me.webApi.UseMiddleware(me.readLanguage)
	me.webApi.UseMiddleware(me.checkAdminMode)
	var appGoals = new(webAppGoals)
	appGoals.init(me.webApiGroup, me.db)

	var appAdmin = new(webAppAdmin)
	appAdmin.init(me.webApiGroup, me.db)
}

func (me *webApp) catchPanic(context huma.Context, next func(huma.Context)) {
	defer func() {
		var exception = recover()
		if exception != nil {
			var webError, isWebError = exception.(webError)
			if isWebError {
				gophers.AssertError(huma.WriteErr(
					me.webApi, context, webError.Status, string(gophers.EncodeJson(webError)),
				))
			} else {
				log.Printf("Error in web function: %v\n%s", exception, debug.Stack())
				gophers.AssertError(huma.WriteErr(
					me.webApi, context, http.StatusInternalServerError, "",
				))
			}
		}
	}()
	next(context)
}

func (webApp) readLanguage(context huma.Context, next func(huma.Context)) {
	var acceptLanguage = ""
	context.EachHeader(func(name, value string) {
		if http.CanonicalHeaderKey(name) == "Accept-Language" {
			acceptLanguage = value
		}
	})
	var languageTag language.Tag = parseLanguageHeader(acceptLanguage)
	context = huma.WithValue(context, webContext.keyLanguage, languageTag)
	next(context)
}

func (webApp) checkAdminMode(context huma.Context, next func(huma.Context)) {
	var adminPasswordCookie, adminPasswordCookieError = huma.ReadCookie(context, "adminPassword")
	var adminPassword = gophers.ReadEnvVar("ADMIN_PASSWORD", "")
	var isAdmin = len(adminPassword) > 0 &&
		adminPasswordCookieError == nil &&
		adminPasswordCookie.Value == adminPassword
	context = huma.WithValue(context, webContext.keyIsAdmin, isAdmin)
	next(context)
}
