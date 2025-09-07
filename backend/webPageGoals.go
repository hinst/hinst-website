package main

import (
	"net/http"

	"github.com/hinst/hinst-website/pages"
)

type webPageGoals struct {
	webAppGoalsBase
	webPath string
}

func (me *webPageGoals) init(db *database, webPath string) []namedWebFunction {
	me.db = db
	me.webPath = webPath

	var fileServer = http.FileServer(http.Dir("./pages/static"))
	var filesPrefix = me.webPath + "/pages/static/"
	http.Handle(filesPrefix, http.StripPrefix(filesPrefix, fileServer))

	return []namedWebFunction{
		{"/pages", me.getHomePage},
	}
}

func (me *webPageGoals) getHomePage(response http.ResponseWriter, request *http.Request) {
	writeHtmlResponse(response, me.getTemplatePage("test home"))
}

func (me *webPageGoals) getTemplatePage(content string) string {
	var page = pages.Template{
		BaseTemplate: me.getBaseTemplate(),
		Content:      content,
	}
	return executeTemplateFile("pages/template.html", page)
}

func (me *webPageGoals) getBaseTemplate() pages.BaseTemplate {
	return pages.BaseTemplate{WebPath: me.webPath}
}
