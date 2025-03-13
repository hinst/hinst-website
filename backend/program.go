package main

import (
	"log"
	"net/http"
)

type program struct {
	savedGoalsPath string
}

func (me *program) init() *program {
	me.savedGoalsPath = "./saved-goals"
	return me
}

func (me *program) runWeb() {
	const netAddress = ":8080"
	var webApp = &webApp{savedGoalsPath: me.savedGoalsPath}
	webApp.init()
	log.Printf("Starting: netAddress=%v, webPath=%v", netAddress, webApp.webPath)
	assertError(http.ListenAndServe(netAddress, nil))
}

func (me *program) runTranslate() {
	(&translator{savedGoalsPath: me.savedGoalsPath}).run()
}
