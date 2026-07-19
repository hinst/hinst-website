package server

import (
	"net/http"
	"time"

	"github.com/hinst/go-gophers"
)

const pagesWebPath = "/pages"

// Server side rendering for personal goals pages.
// This code is used to generate static files to be displayed on a hosting service without backend API.
type webPageGoals struct {
	webAppGoalsBase
	renderer *goalRenderer
	webPath  string
}

func (me *webPageGoals) init(db *database, webPath string) []namedWebFunction {
	me.db = db
	me.webPath = webPath
	me.renderer = newGoalRenderer(db, webPath)

	var fileServer = http.FileServer(http.Dir("./pages/static"))
	var filesPrefix = me.webPath + "/static/"
	http.Handle(filesPrefix, http.StripPrefix(filesPrefix, fileServer))

	return []namedWebFunction{
		{"", me.getHomePage},
		{"/personal-goals/{id}", me.getGoalPage},
		{"/personal-goals/{id}/{post}", me.getGoalPostPage},
		{"/personal-goals/image/{id}/{post}/{index}", me.getGoalPostImage},
	}
}

func (me *webPageGoals) getHomePage(response http.ResponseWriter, request *http.Request) {
	var req = WebRequest{
		Language:      getWebLanguage(request),
		WebPath:       me.inputWebPath(request.URL.Query().Get("webPath"), me.webPath),
		StaticPath:    me.inputWebPath(request.URL.Query().Get("staticPath"), me.webPath),
		JpegExtension: request.URL.Query().Get("jpegExtension"),
		HtmlExtension: request.URL.Query().Get("htmlExtension"),
	}
	var html = me.renderer.renderHomePage(req)
	writeHtmlResponse(response, html)
}

func (me *webPageGoals) getGoalPage(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.PathValue("id"))
	var req = WebRequest{
		Language:      getWebLanguage(request),
		WebPath:       me.inputWebPath(request.URL.Query().Get("webPath"), me.webPath),
		StaticPath:    me.inputWebPath(request.URL.Query().Get("staticPath"), me.webPath),
		JpegExtension: request.URL.Query().Get("jpegExtension"),
		HtmlExtension: request.URL.Query().Get("htmlExtension"),
	}
	var html = me.renderer.renderGoalPage(req, goalId)
	writeHtmlResponse(response, html)
}

func (me *webPageGoals) getGoalPostPage(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.PathValue("id"))
	var dateTime = me.inputValidPostDateTime(request.PathValue("post"))

	var req = WebRequest{
		Language:      getWebLanguage(request),
		WebPath:       me.inputWebPath(request.URL.Query().Get("webPath"), me.webPath),
		StaticPath:    me.inputWebPath(request.URL.Query().Get("staticPath"), me.webPath),
		JpegExtension: request.URL.Query().Get("jpegExtension"),
		HtmlExtension: request.URL.Query().Get("htmlExtension"),
	}
	var html = me.renderer.renderGoalPostPage(req, goalId, dateTime)
	writeHtmlResponse(response, html)
}

func (me *webPageGoals) getGoalPostImage(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.PathValue("id"))
	var postDateTime = me.inputValidPostDateTime(request.PathValue("post"))
	var index = inputValidWebInteger(request.PathValue("index"))
	var image = me.db.getGoalPostImage(goalId, postDateTime, index)
	if image == nil {
		panic(webError{"Image not found", http.StatusNotFound})
	}
	gophers.SetCacheAge(response, time.Hour)
	response.Header().Set(gophers.ContentTypeHeader, image.ContentType)
	var _, _ = response.Write(image.File)
}

func (me *webPageGoals) inputWebPath(text string, defaultText string) string {
	if text == "" {
		return defaultText
	}
	if text == "/" {
		return ""
	}
	return text
}
