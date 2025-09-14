package server

import (
	"html/template"
	"net/http"

	"github.com/hinst/hinst-website/server/page_data"
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
	var data = page_data.GoalList{Base: me.getBaseTemplate()}
	for _, goal := range goals {
		var item page_data.GoalCard
		var metaInfo = goalInfo{}.findByTitle(personalGoalInfos, goal.Title)
		item.Id = goal.Id
		item.Title = goal.Title
		item.Image = metaInfo.coverImage
		if language != supportedLanguages[0] {
			item.Title = metaInfo.englishTitle
		}
		data.Goals = append(data.Goals, item)
	}
	var content = executeTemplateFile("pages/html/templates/goalList.html", data)
	writeHtmlResponse(response, me.getTemplatePage("My Personal Goals", content))
}

func (me *webPageGoals) getGoalPage(response http.ResponseWriter, request *http.Request) {
	var requestedLanguage = getWebLanguage(request)
	var goalId = me.inputValidGoalId(request.URL.Query().Get("id"))
	var goalPosts = me.db.getGoalPosts(goalId, false, requestedLanguage)
	if goalPosts == nil {
		var errorMessage = "Cannot find goalId=" + getStringFromInt64(goalId)
		panic(webError{errorMessage, http.StatusNotFound})
	}
}

func (me *webPageGoals) getTemplatePage(title string, content string) string {
	var headerData = page_data.ContentTemplate{Base: me.getBaseTemplate(), Title: title}
	var page = page_data.ContentTemplate{
		Base:    me.getBaseTemplate(),
		Title:   title,
		Header:  template.HTML(executeTemplateFile("pages/html/templates/header.html", headerData)),
		Content: template.HTML(content),
	}
	return executeTemplateFile("pages/html/templates/template.html", page)
}

func (me *webPageGoals) getBaseTemplate() page_data.Base {
	return page_data.Base{WebPath: me.webPath}
}
