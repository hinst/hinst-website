package main

import (
	"html/template"
	"net/http"
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
	var goals = me.db.getGoals()
	var data = GoalListTemplate{BaseTemplate: me.getBaseTemplate(), Goals: goals}
	var content = executeTemplateFile("pages/goalList.html", data)
	writeHtmlResponse(response, me.getTemplatePage("My Personal Goals", content))
}

func (me *webPageGoals) getTemplatePage(title string, content string) string {
	var headerData = ContentTemplate{BaseTemplate: me.getBaseTemplate(), Title: title}
	var page = ContentTemplate{
		BaseTemplate: me.getBaseTemplate(),
		Title:        title,
		Header:       template.HTML(executeTemplateFile("pages/header.html", headerData)),
		Content:      template.HTML(content),
	}
	return executeTemplateFile("pages/template.html", page)
}

func (me *webPageGoals) getBaseTemplate() BaseTemplate {
	return BaseTemplate{WebPath: me.webPath}
}
