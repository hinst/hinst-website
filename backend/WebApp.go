package main

import (
	"net/http"
	"os"
)

type webApp struct {
	path string
}

func (me *webApp) start() {
	http.HandleFunc(me.path+"/api/get-goals", me.getGoals)
}

func (me *webApp) getGoals(writer http.ResponseWriter, request *http.Request) {
	var files = assertResultError(os.ReadDir("./saved-goals"))
	
}
