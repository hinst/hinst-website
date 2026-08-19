package server

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/hinst/go-gophers"
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
	me.webApi.UseMiddleware(func(context huma.Context, next func(huma.Context)) {
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
	})
	var appGoals = new(webAppGoals)
	appGoals.init(me.webApiGroup, me.db)

	var appAdmin = new(webAppAdmin)
	appAdmin.init(me.db)
}
