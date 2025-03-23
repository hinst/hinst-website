package main

import (
	"log"
	"net/http"
)

type program struct {
	netAddress       string
	webFilesPath     string
	savedGoalsPath   string
	allowOrigin      string
	translatorApiUrl string
}

var programTemplate = program{
	netAddress:       ":8080",
	webFilesPath:     "./www",
	savedGoalsPath:   "./saved-goals",
	allowOrigin:      "http://localhost:1234",
	translatorApiUrl: "http://localhost:1235",
}

func (me *program) create() *program {
	*me = programTemplate
	return me
}

func (me *program) runWeb() {
	var webApp = &webApp{
		savedGoalsPath: me.savedGoalsPath,
		allowOrigin:    me.allowOrigin,
	}
	webApp.init()
	var fileServer = http.FileServer(http.Dir(me.webFilesPath))
	http.Handle(webApp.webPath+"/", http.StripPrefix(webApp.webPath+"/", fileServer))
	log.Printf("Starting: netAddress=%v, webPath=%v, filesPath=%v",
		me.netAddress, webApp.webPath, me.webFilesPath)
	assertError(http.ListenAndServe(me.netAddress, nil))
}

func (me *program) runTranslate() {
	var translator = translatorPresets
	if me.translatorApiUrl != "" {
		translator.apiUrl = me.translatorApiUrl + "/v1/chat/completions"
	}
	translator.savedGoalsPath = me.savedGoalsPath
	translator.run()
}
