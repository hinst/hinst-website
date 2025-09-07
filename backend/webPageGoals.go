package main

import (
	"net/http"

	"github.com/hinst/hinst-website/pages"
)

type webPageGoals struct {
	webAppGoalsBase
}

func (me *webPageGoals) init(db *database) []namedWebFunction {
	me.db = db
	return []namedWebFunction{
		{"/static", me.getHomePage},
	}
}

func (me *webPageGoals) getHomePage(response http.ResponseWriter, request *http.Request) {
	writeHtmlResponse(response, me.getTemplatePage("test home"))
}

func (me *webPageGoals) getTemplatePage(content string) string {
	var page = pages.Template{Content: content}
	return executeTemplateFile("pages/template.html", page)

}
