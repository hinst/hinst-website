package main

import (
	"net/http"
	"os"
)

type webApp struct {
	path string
}

func (me *webApp) start() {
	http.HandleFunc(me.path+"/api/goals", me.wrap(me.getGoals))
}

func (me *webApp) wrap(function func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost:1234")
		function(response, request)
	}
}

func (me *webApp) getGoals(response http.ResponseWriter, request *http.Request) {
	var files = assertResultError(os.ReadDir("./saved-goals"))
	var headers = make([]*goalHeader, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			var headerFilePath = "./saved-goals/" + file.Name() + "/_header.json"
			var header = readJsonFile(headerFilePath, &goalHeader{})
			headers = append(headers, header)
		}
	}
	response.Write(encodeJson(headers))
}
