package server

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
		{pagesWebPath, me.getHomePage},
	}
}

func (me *webPageGoals) getHomePage(response http.ResponseWriter, request *http.Request) {
	var language = getWebLanguage(request)
	var goals = me.db.getGoals()
	var data = goalListTemplate{BaseTemplate: me.getBaseTemplate()}
	for _, goal := range goals {
		var item goalCardTemplate
		var metaInfo = goalInfo{}.findByTitle(personalGoalInfos, goal.Title)
		item.Id = goal.Id
		item.Title = goal.Title
		item.Image = metaInfo.coverImage
		if language != supportedLanguages[0] {
			item.Title = metaInfo.englishTitle
		}
		data.Goals = append(data.Goals, item)
	}
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
