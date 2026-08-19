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
	var appGoals = new(webAppGoals)
	me.addFunctions(me.webPath, appGoals.init(me.webApiGroup, me.db))

	var appAdmin = new(webAppAdmin)
	me.addFunctions(me.webPath, appAdmin.init(me.db))
}

func (me *webApp) addFunctions(path string, functions []namedWebFunction) {
	for _, namedWebFunction := range functions {
		var url = path + namedWebFunction.Name
		log.Printf("Adding web function: %v", url)
		http.HandleFunc(url, me.wrap(namedWebFunction.Function))
	}
}

func (me *webApp) wrap(function gophers.WebFunction) gophers.WebFunction {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", me.allowOrigin())
		defer func() {
			var exception = recover()
			if exception != nil {
				var webError, isWebError = exception.(webError)
				if isWebError {
					response.WriteHeader(webError.Status)
					var _, _ = response.Write(gophers.EncodeJson(webError))
				} else {
					response.WriteHeader(http.StatusInternalServerError)
					log.Printf("Error in web function: %v\n%s", exception, debug.Stack())
				}
			}
		}()
		function(response, request)
	}
}

func (webApp) allowOrigin() string {
	return gophers.ReadEnvVar("ALLOW_ORIGIN", "")
}
