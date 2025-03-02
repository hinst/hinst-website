package main

import (
	"net/http"
	"os"
)

type webApp struct {
	path string
}

func (me *webApp) start() {
	http.HandleFunc(me.path+"/api/goals", me.getGoals)
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
	response.WriteHeader(http.StatusOK)
	response.Write(encodeJson(headers))
}
