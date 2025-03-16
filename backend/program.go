package main

import (
	"log"
	"net/http"
)

type program struct {
	webFilesPath   string
	savedGoalsPath string
}

func (me *program) init() *program {
	me.webFilesPath = "./www"
	me.savedGoalsPath = "./saved-goals"
	return me
}

func (me *program) runWeb() {
	const netAddress = ":8080"
	var webApp = &webApp{savedGoalsPath: me.savedGoalsPath}
	webApp.init()
	http.Handle("/", http.FileServer(http.Dir(me.webFilesPath)))
	log.Printf("Starting: netAddress=%v, webPath=%v", netAddress, webApp.webPath)
	assertError(http.ListenAndServe(netAddress, nil))
}

func (me *program) runTranslate() {
	var translator = translatorPresets
	translator.savedGoalsPath = me.savedGoalsPath
	translator.run()
}
