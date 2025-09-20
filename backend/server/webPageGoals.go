package server

import (
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/hinst/hinst-website/server/page_data"
	"golang.org/x/text/language"
)

type webPageGoals struct {
	webAppGoalsBase
	webPath         string
	elementIdLocker sync.Mutex
	elementId       int64
}

func (me *webPageGoals) init(db *database, webPath string) []namedWebFunction {
	me.db = db
	me.webPath = webPath

	var fileServer = http.FileServer(http.Dir("./pages/static"))
	var filesPrefix = me.webPath + "/static/"
	http.Handle(filesPrefix, http.StripPrefix(filesPrefix, fileServer))

	return []namedWebFunction{
		{"", me.getHomePage},
		{"/personal-goals/{id}", me.getGoalPage},
		{"/personal-goals/{id}/{post}", me.getGoalPostPage},
		{"/personal-goals/image", me.getGoalPostImage},
	}
}

func (me *webPageGoals) getHomePage(response http.ResponseWriter, request *http.Request) {
	var language = getWebLanguage(request)
	var goals = me.db.getGoals()
	var data = page_data.GoalList{Base: me.getBaseTemplate(request)}
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
	writeHtmlResponse(response, me.wrapTemplatePage(request, "My Personal Goals", content))
}

func (me *webPageGoals) getGoalPage(response http.ResponseWriter, request *http.Request) {
	var requestedLanguage = getWebLanguage(request)
	var goalId = me.inputValidGoalId(request.PathValue("id"))
	var goalPostRecords = me.db.getGoalPosts(goalId, false, requestedLanguage)
	if goalPostRecords == nil {
		var errorMessage = "Cannot find goal with id=" + getStringFromInt64(goalId)
		panic(webError{errorMessage, http.StatusNotFound})
	}

	var goalPosts []page_data.GoalPostItem
	for _, post := range goalPostRecords {
		if post.Title == nil {
			continue
		}
		var item page_data.GoalPostItem
		item.Title = *post.Title
		item.DateTime = post.DateTime
		item.Day = time.Unix(post.DateTime, 0).UTC().Day()
		goalPosts = append(goalPosts, item)
	}

	var data = page_data.GoalPosts{Base: me.getBaseTemplate(request)}
	data.GoalId = goalId
	data.Load(goalPosts)

	var content = executeTemplateFile("pages/html/templates/goalPosts.html", data)
	writeHtmlResponse(response, me.wrapTemplatePage(request, "Goal diary", content))
}

func (me *webPageGoals) getGoalPostPage(response http.ResponseWriter, request *http.Request) {
	var requestedLanguage = getWebLanguage(request)
	var goalId = me.inputValidGoalId(request.PathValue("id"))
	var dateTime = me.inputValidPostDateTime(request.PathValue("post"))

	var goalRecord = me.db.getGoal(goalId)
	if goalRecord == nil {
		var errorMessage = "Cannot find goal with id=" + getStringFromInt64(goalId)
		panic(webError{errorMessage, http.StatusNotFound})
	}

	var goalPostRecord = me.db.getGoalPost(goalId, dateTime)
	if goalPostRecord == nil {
		var errorMessage = "Cannot find goal post with id=" + getStringFromInt64(goalId) +
			" and dateTime=" + dateTime.UTC().Format(time.DateTime)
		panic(webError{errorMessage, http.StatusNotFound})
	}
	var text = goalPostRecord.getTranslatedText(requestedLanguage)
	var data = page_data.GoalPost{
		Base:     me.getBaseTemplate(request),
		GoalId:   goalId,
		DateTime: dateTime.Unix(),
		Text:     template.HTML(text),
	}

	var imageCount = me.db.getGoalPostImageCount(goalId, dateTime)
	for i := range imageCount {
		data.Images = append(data.Images, i)
	}

	var pageTitle = me.getTranslatedTitle(goalRecord.Title, requestedLanguage) + " • " +
		dateTime.UTC().Format("2006-01-02")
	var content = executeTemplateFile("pages/html/templates/goalPost.html", data)
	var goalTitle = me.getTranslatedTitle(pageTitle, requestedLanguage)
	writeHtmlResponse(response, me.wrapTemplatePage(request, goalTitle, content))
}

func (me *webPageGoals) getTranslatedTitle(title string, language language.Tag) string {
	if language == supportedLanguages[0] {
		return title
	}
	var metaInfo = goalInfo{}.findByTitle(personalGoalInfos, title)
	if metaInfo == nil {
		return title
	}
	return metaInfo.englishTitle
}

func (me *webPageGoals) wrapTemplatePage(request *http.Request, pageTitle string, content string) string {
	var headerData = page_data.Content{Base: me.getBaseTemplate(request), Title: pageTitle}
	var page = page_data.Content{
		Base:    me.getBaseTemplate(request),
		Title:   pageTitle,
		Header:  template.HTML(executeTemplateFile("pages/html/templates/header.html", headerData)),
		Content: template.HTML(content),
	}
	return executeTemplateFile("pages/html/templates/template.html", page)
}

func (me *webPageGoals) getGoalPostImage(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("id"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var index = requireRequestQueryInt(request, "index")
	var image = me.db.getGoalPostImage(goalId, postDateTime, index)
	if image == nil {
		panic(webError{"Image not found", http.StatusNotFound})
	}
	setCacheAge(response, time.Hour)
	response.Header().Set("Content-Type", image.contentType)
	response.Write(image.file)
}

func (me *webPageGoals) getBaseTemplate(request *http.Request) page_data.Base {
	var webPath = me.webPath
	var customWebPath = request.URL.Query().Get("webPath")
	if customWebPath != "" {
		webPath = customWebPath
		if webPath == "/" {
			webPath = ""
		}
	}

	var apiPath = webPath
	var customApiPath = request.URL.Query().Get("apiPath")
	if customApiPath != "" {
		apiPath = customApiPath
		if apiPath == "/" {
			apiPath = ""
		}
	}

	return page_data.Base{
		Id:          me.advanceElementId(),
		WebPath:     webPath,
		ApiPath:     apiPath,
		SettingsSvg: template.HTML(readTextFile("pages/static/images/settings.svg")),
	}
}

func (me *webPageGoals) advanceElementId() int64 {
	me.elementIdLocker.Lock()
	defer me.elementIdLocker.Unlock()
	me.elementId++
	return me.elementId
}
