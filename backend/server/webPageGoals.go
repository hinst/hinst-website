package server

import (
	"html/template"
	"net/http"
	"slices"
	"time"

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
		{pagesWebPath + "/personal-goals", me.getGoalPage},
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
	var data = page_data.GoalPosts{Base: me.getBaseTemplate()}
	data.GoalId = goalId
	for _, post := range goalPosts {
		if post.Title == nil {
			continue
		}
		var item page_data.GoalPost
		item.Title = *post.Title
		item.DateTime = post.DateTime
		item.DateText = time.Unix(post.DateTime, 0).Format("2006-01-02")
		data.Posts = append(data.Posts, item)
	}
	slices.SortFunc(data.Posts, func(a, b page_data.GoalPost) int {
		if a.DateTime < b.DateTime {
			return 1
		} else if a.DateTime > b.DateTime {
			return -1
		} else {
			return 0
		}
	})
	var content = executeTemplateFile("pages/html/templates/goalPosts.html", data)
	writeHtmlResponse(response, me.getTemplatePage("Goal diary", content))
}

func (me *webPageGoals) getTemplatePage(title string, content string) string {
	var headerData = page_data.Content{Base: me.getBaseTemplate(), Title: title}
	var page = page_data.Content{
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
