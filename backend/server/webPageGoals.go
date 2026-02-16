package server

import (
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/hinst/hinst-website/server/page_data"
	"golang.org/x/text/language"
)

// Server side rendering for personal goals pages.
// This code is used to generate static files to be displayed on a hosting service without backend API.
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
		{"/personal-goals/image/{id}/{post}/{index}", me.getGoalPostImage},
	}
}

func (me *webPageGoals) getHomePage(response http.ResponseWriter, request *http.Request) {
	var requestedLanguage = getWebLanguage(request)
	var goalRecords = me.db.getGoals()
	var data = page_data.GoalList{Base: me.getBaseTemplate(request)}
	for _, goalRecord := range goalRecords {
		var item page_data.GoalCard
		var metaInfo = goalInfo{}.findByTitle(personalGoalInfos, goalRecord.Title)
		item.Id = goalRecord.Id
		item.Title = goalRecord.Title
		item.Image = metaInfo.coverImage
		if requestedLanguage != supportedLanguages[0] {
			item.Title = metaInfo.englishTitle
		}
		data.Goals = append(data.Goals, item)
	}
	var content = executeTemplateFile("pages/html/templates/goalList.html", data)
	writeHtmlResponse(response, me.wrapTemplatePage(request, page_data.Content{
		LanguageTag: requestedLanguage.String(),
		Title:       "My Personal Goals",
		Content:     template.HTML(content),
	}))
}

func (me *webPageGoals) getGoalPage(response http.ResponseWriter, request *http.Request) {
	var requestedLanguage = getWebLanguage(request)
	var goalId = me.inputValidGoalId(request.PathValue("id"))
	var goalRecord = me.db.getGoal(goalId)
	if goalRecord == nil {
		var errorMessage = "Cannot find goal with id=" + getStringFromInt64(goalId)
		panic(webError{errorMessage, http.StatusNotFound})
	}
	var goalPostRecords = me.db.getGoalPosts(goalId, false, requestedLanguage)

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

	var metaInfo = goalInfo{}.findByTitle(personalGoalInfos, goalRecord.Title)
	var goalTitle = goalRecord.Title
	if requestedLanguage != supportedLanguages[0] {
		goalTitle = metaInfo.englishTitle
	}
	var content = executeTemplateFile("pages/html/templates/goalPosts.html", data)
	writeHtmlResponse(response, me.wrapTemplatePage(request, page_data.Content{
		LanguageTag: requestedLanguage.String(),
		Title:       "Goal diary: " + goalTitle,
		Content:     template.HTML(content),
	}))
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

	var goalTitle = me.getTranslatedTitle(goalRecord.Title, requestedLanguage)
	var pageTitle = goalTitle + " • " +
		dateTime.UTC().Format("2006-01-02")
	var pageDescription = goalTitle + " - " +
		dateTime.UTC().Format("2006-01-02") + " - " +
		goalPostRecord.getTranslatedTitle(requestedLanguage)
	var content = executeTemplateFile("pages/html/templates/goalPost.html", data)
	writeHtmlResponse(response, me.wrapTemplatePage(request, page_data.Content{
		LanguageTag: requestedLanguage.String(),
		Title:       pageTitle,
		Description: pageDescription,
		Content:     template.HTML(content),
	}))
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

func (me *webPageGoals) wrapTemplatePage(request *http.Request, content page_data.Content) string {
	if content.Description == "" {
		content.Description = content.Title
	}
	var headerContent = content
	headerContent.Base = me.getBaseTemplate(request)
	var pageContent = content
	pageContent.Base = me.getBaseTemplate(request)
	pageContent.Header = template.HTML(executeTemplateFile("pages/html/templates/header.html", headerContent))
	return executeTemplateFile("pages/html/templates/template.html", pageContent)
}

func (me *webPageGoals) getGoalPostImage(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.PathValue("id"))
	var postDateTime = me.inputValidPostDateTime(request.PathValue("post"))
	var index = inputValidWebInteger(request.PathValue("index"))
	var image = me.db.getGoalPostImage(goalId, postDateTime, index)
	if image == nil {
		panic(webError{"Image not found", http.StatusNotFound})
	}
	setCacheAge(response, time.Hour)
	response.Header().Set("Content-Type", image.contentType)
	var _, _ = response.Write(image.file)
}

func (me *webPageGoals) getBaseTemplate(request *http.Request) page_data.Base {
	return page_data.Base{
		Id:            me.advanceElementId(),
		WebPath:       me.inputWebPath(request.URL.Query().Get("webPath"), me.webPath),
		StaticPath:    me.inputWebPath(request.URL.Query().Get("staticPath"), me.webPath),
		JpegExtension: request.URL.Query().Get("jpegExtension"),
		HtmlExtension: request.URL.Query().Get("htmlExtension"),
		SettingsSvg:   template.HTML(readTextFile("pages/static/images/settings.svg")),
		MenuSvg:       template.HTML(readTextFile("pages/static/images/menu.svg")),
	}
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

func (me *webPageGoals) advanceElementId() int64 {
	me.elementIdLocker.Lock()
	defer me.elementIdLocker.Unlock()
	me.elementId++
	return me.elementId
}
