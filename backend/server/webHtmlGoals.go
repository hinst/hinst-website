package server

import (
	"html/template"
	"sync/atomic"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/page_data"
	"golang.org/x/text/language"
)

type webHtmlGoals struct {
	db        *database
	elementId atomic.Int64
}

func (me *webHtmlGoals) renderHomePage(webPath string, lang language.Tag) string {
	var goalRecords = me.db.getGoals()
	var langPath = webStaticGoals{}.getLanguagePath(lang)
	var data = page_data.GoalList{Base: me.getBaseTemplate(webPath, langPath)}
	for _, goalRecord := range goalRecords {
		var item page_data.GoalCard
		item.Id = goalRecord.Id
		var imageDataUrl = getUrlBase64(goalRecord.ImageContentType, goalRecord.ImageData)
		item.Image = template.URL(imageDataUrl)
		item.Title = goalRecord.GetTranslatedTitle(lang)
		data.Goals = append(data.Goals, item)
	}
	var content = executeTemplateFile("pages/html/templates/goalList.html", data)
	return me.wrapTemplatePage(webPath, langPath, page_data.Content{
		LanguageTag: lang.String(),
		Title:       "My Personal Goals",
		Content:     template.HTML(content),
	})
}

func (me *webHtmlGoals) renderGoalPage(webPath string, lang language.Tag, goalId int64) string {
	var goalRecord = me.db.getGoal(goalId)
	gophers.AssertCondition(goalRecord != nil, func() string { return "Cannot find goal with id=" + gophers.GetStringFromInt64(goalId) })
	var goalPostRecords = me.db.getGoalPosts(goalId, false, lang)
	var langPath = webStaticGoals{}.getLanguagePath(lang)

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

	var data = page_data.GoalPosts{Base: me.getBaseTemplate(webPath, langPath)}
	data.GoalId = goalId
	data.Load(goalPosts)

	var content = executeTemplateFile("pages/html/templates/goalPosts.html", data)
	return me.wrapTemplatePage(webPath, langPath, page_data.Content{
		LanguageTag: lang.String(),
		Title:       "Goal diary: " + goalRecord.GetTranslatedTitle(lang),
		Content:     template.HTML(content),
	})
}

func (me *webHtmlGoals) renderGoalPostPage(webPath string, lang language.Tag, goalId int64, dateTime time.Time) string {
	var goalRecord = me.db.getGoal(goalId)
	gophers.AssertCondition(goalRecord != nil, func() string { return "Cannot find goal with id=" + gophers.GetStringFromInt64(goalId) })
	var langPath = webStaticGoals{}.getLanguagePath(lang)

	var goalPostRecord = me.db.getGoalPost(goalId, dateTime)
	gophers.AssertCondition(goalPostRecord != nil, func() string {
		return "Cannot find goal post with id=" + gophers.GetStringFromInt64(goalId) +
			" and dateTime=" + dateTime.UTC().Format(time.DateTime)
	})

	var text = goalPostRecord.GetTranslatedText(lang)
	var data = page_data.GoalPost{
		Base:         me.getBaseTemplate(webPath, langPath),
		GoalId:       goalId,
		DateTime:     dateTime.Unix(),
		Text:         template.HTML(convertMarkdownToHtml(text)),
		LanguageName: base.GetLanguageName(lang),
	}
	if lang != base.SupportedLanguages[0] {
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

	var goalTitle = goalRecord.GetTranslatedTitle(lang)
	var pageTitle = goalTitle + " • " +
		dateTime.UTC().Format("2006-01-02")
	var pageDescription = goalTitle + " - " +
		dateTime.UTC().Format("2006-01-02") + " - " +
		goalPostRecord.GetTranslatedTitle(lang)
	var content = executeTemplateFile("pages/html/templates/goalPost.html", data)
	return me.wrapTemplatePage(webPath, langPath, page_data.Content{
		LanguageTag: lang.String(),
		Title:       pageTitle,
		Description: pageDescription,
		Content:     template.HTML(content),
	})
}

func (me *webHtmlGoals) wrapTemplatePage(webPath string, langPath string, content page_data.Content) string {
	if content.Description == "" {
		content.Description = content.Title
	}
	headerContent := content
	headerContent.Base = me.getBaseTemplate(webPath, langPath)
	htmlHeader := executeTemplateFile("pages/html/templates/header.html", headerContent)

	pageContent := content
	pageContent.Base = me.getBaseTemplate(webPath, langPath)
	pageContent.Header = template.HTML(htmlHeader)
	return executeTemplateFile("pages/html/templates/template.html", pageContent)
}

func (me *webHtmlGoals) getBaseTemplate(webPath string, langPath string) page_data.Base {
	return page_data.Base{
		Id:          me.elementId.Add(1),
		WebPath:     webPath,
		LangPath:    langPath,
		SettingsSvg: template.HTML(gophers.ReadTextFile("pages/static/images/settings.svg")),
		MenuSvg:     template.HTML(gophers.ReadTextFile("pages/static/images/menu.svg")),
		InfoSvg:     template.HTML(gophers.ReadTextFile("pages/static/images/info.svg")),
	}
}
