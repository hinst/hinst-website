package server

import (
	"html/template"
	"sync/atomic"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/page_data"
)

type goalRenderer struct {
	db        *database
	elementId atomic.Int64
}

func (me *goalRenderer) renderHomePage(req WebRequest) string {
	var goalRecords = me.db.getGoals()
	var data = page_data.GoalList{Base: me.getBaseTemplate(req)}
	for _, goalRecord := range goalRecords {
		var item page_data.GoalCard
		item.Id = goalRecord.Id
		var imageDataUrl = getUrlBase64(goalRecord.ImageContentType, goalRecord.ImageData)
		item.Image = template.URL(imageDataUrl)
		item.Title = goalRecord.GetTranslatedTitle(req.Language)
		data.Goals = append(data.Goals, item)
	}
	var content = executeTemplateFile("pages/html/templates/goalList.html", data)
	return me.wrapTemplatePage(req, page_data.Content{
		LanguageTag: req.Language.String(),
		Title:       "My Personal Goals",
		Content:     template.HTML(content),
	})
}

func (me *goalRenderer) renderGoalPage(req WebRequest, goalId int64) string {
	var goalRecord = me.db.getGoal(goalId)
	gophers.AssertCondition(goalRecord != nil, func() string { return "Cannot find goal with id=" + gophers.GetStringFromInt64(goalId) })

	var goalPostRecords = me.db.getGoalPosts(goalId, false, req.Language)

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

	var data = page_data.GoalPosts{Base: me.getBaseTemplate(req)}
	data.GoalId = goalId
	data.Load(goalPosts)

	var content = executeTemplateFile("pages/html/templates/goalPosts.html", data)
	return me.wrapTemplatePage(req, page_data.Content{
		LanguageTag: req.Language.String(),
		Title:       "Goal diary: " + goalRecord.GetTranslatedTitle(req.Language),
		Content:     template.HTML(content),
	})
}

func (me *goalRenderer) renderGoalPostPage(req WebRequest, goalId int64, dateTime time.Time) string {
	var goalRecord = me.db.getGoal(goalId)
	gophers.AssertCondition(goalRecord != nil, func() string { return "Cannot find goal with id=" + gophers.GetStringFromInt64(goalId) })

	var goalPostRecord = me.db.getGoalPost(goalId, dateTime)
	gophers.AssertCondition(goalPostRecord != nil, func() string {
		return "Cannot find goal post with id=" + gophers.GetStringFromInt64(goalId) +
			" and dateTime=" + dateTime.UTC().Format(time.DateTime)
	})

	var text = goalPostRecord.GetTranslatedText(req.Language)
	var data = page_data.GoalPost{
		Base:         me.getBaseTemplate(req),
		GoalId:       goalId,
		DateTime:     dateTime.Unix(),
		Text:         template.HTML(convertMarkdownToHtml(text)),
		LanguageName: base.GetLanguageName(req.Language),
	}
	if req.Language != base.SupportedLanguages[0] {
		if text == "" {
			data.IsTranslationPending = true
			data.Text = template.HTML(convertMarkdownToHtml(goalPostRecord.Text))
		} else {
			data.IsAutoTranslated = true
		}
	}

	var imageCount = me.db.getGoalPostImageCount(goalId, dateTime)
	for i := range imageCount {
		data.Images = append(data.Images, i)
	}

	var goalTitle = goalRecord.GetTranslatedTitle(req.Language)
	var pageTitle = goalTitle + " • " +
		dateTime.UTC().Format("2006-01-02")
	var pageDescription = goalTitle + " - " +
		dateTime.UTC().Format("2006-01-02") + " - " +
		goalPostRecord.GetTranslatedTitle(req.Language)
	var content = executeTemplateFile("pages/html/templates/goalPost.html", data)
	return me.wrapTemplatePage(req, page_data.Content{
		LanguageTag: req.Language.String(),
		Title:       pageTitle,
		Description: pageDescription,
		Content:     template.HTML(content),
	})
}

func (me *goalRenderer) wrapTemplatePage(req WebRequest, content page_data.Content) string {
	if content.Description == "" {
		content.Description = content.Title
	}
	headerContent := content
	headerContent.Base = me.getBaseTemplate(req)
	htmlHeader := executeTemplateFile("pages/html/templates/header.html", headerContent)

	pageContent := content
	pageContent.Base = me.getBaseTemplate(req)
	pageContent.Header = template.HTML(htmlHeader)
	return executeTemplateFile("pages/html/templates/template.html", pageContent)
}

func (me *goalRenderer) getBaseTemplate(req WebRequest) page_data.Base {
	var webPath = req.WebPath
	var staticPath = req.StaticPath
	if staticPath == "" {
		staticPath = ""
	}

	return page_data.Base{
		Id:          me.elementId.Add(1),
		WebPath:     webPath,
		StaticPath:  staticPath,
		SettingsSvg: template.HTML(gophers.ReadTextFile("pages/static/images/settings.svg")),
		MenuSvg:     template.HTML(gophers.ReadTextFile("pages/static/images/menu.svg")),
		InfoSvg:     template.HTML(gophers.ReadTextFile("pages/static/images/info.svg")),
	}
}
