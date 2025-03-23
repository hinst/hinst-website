package main

import (
	"log"
	"net/http"
)

type program struct {
	webFilesPath     string
	savedGoalsPath   string
	allowOrigin      string
	translatorApiUrl string
}

func (me *program) init() *program {
	me.webFilesPath = "./www"
	me.savedGoalsPath = "./saved-goals"
	return me
}

func (me *program) runWeb() {
	const netAddress = ":8080"
	var webApp = &webApp{
		savedGoalsPath: me.savedGoalsPath,
		allowOrigin:    me.allowOrigin,
	}
	webApp.init()
	var fileServer = http.FileServer(http.Dir(me.webFilesPath))
	http.Handle(webApp.webPath+"/", http.StripPrefix(webApp.webPath+"/", fileServer))
	log.Printf("Starting: netAddress=%v, webPath=%v, filesPath=%v", netAddress, webApp.webPath, me.webFilesPath)
	assertError(http.ListenAndServe(netAddress, nil))
}

func (me *program) runTranslate() {
	var translator = translatorPresets
	if me.translatorApiUrl != "" {
		translator.apiUrl = me.translatorApiUrl + "/v1/chat/completions"
	}
	translator.savedGoalsPath = me.savedGoalsPath
	translator.run()
}
