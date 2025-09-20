package server

import (
	"net/url"
	"os"

	"golang.org/x/text/language"
)

type webStaticGoals struct {
	folder string
	url    string
	db     *database
}

func (me *webStaticGoals) init(url string, db *database) {
	me.folder = "static"
	me.url = url
	me.db = db
}

func (me *webStaticGoals) run() {
	assertError(os.RemoveAll(me.folder))
	assertError(os.MkdirAll(me.folder, os.ModePerm))
	os.CopyFS(me.folder+"/static", os.DirFS("pages/static"))
	for _, lang := range supportedLanguages {
		me.generate(lang)
	}
}

func (me *webStaticGoals) generate(lang language.Tag) {
	var path = me.folder + me.getLanguagePath(lang)
	assertError(os.MkdirAll(path, os.ModePerm))
	var homeUrl = me.url + "/pages?" + me.getPathQuery(lang)
	var homePageText = readTextFromUrl(homeUrl)
	writeTextFile(path+"/index.html", homePageText)

	var goals = me.db.getGoals()
	var goalsPath = path + "/personal-goals"
	assertError(os.MkdirAll(goalsPath, os.ModePerm))
	for _, goal := range goals {
		me.generateGoal(lang, goalsPath, goal)
	}
}

func (me *webStaticGoals) generateGoal(lang language.Tag, goalsPath string, goal goalRecord) {
	var goalId = goal.Id
	var url = me.url + pagesWebPath + "/personal-goals/" + getStringFromInt64(goalId) + "?" +
		me.getPathQuery(lang)
	var goalPageText = readTextFromUrl(url)
	writeTextFile(goalsPath+"/"+getStringFromInt64(goalId)+".html", goalPageText)

	var path = goalsPath + "/" + getStringFromInt64(goalId)
	assertError(os.MkdirAll(path, os.ModePerm))
	var posts = me.db.getGoalPosts(goalId, false, lang)
	for _, post := range posts {
		me.generateGoalPost(lang, goalsPath, goalId, post.DateTime)
	}
}

func (me *webStaticGoals) generateGoalPost(lang language.Tag, goalsPath string, goalId int64, postDateTime int64) {
	var url = me.url + pagesWebPath + "/personal-goals/" + getStringFromInt64(goalId) + "/" +
		getStringFromInt64(postDateTime) + "?" + me.getPathQuery(lang)
	var postPageText = readTextFromUrl(url)
	var path = goalsPath + "/" + getStringFromInt64(goalId)
	writeTextFile(path+"/"+getStringFromInt64(postDateTime)+".html", postPageText)
}

func (me *webStaticGoals) getLanguagePath(tag language.Tag) string {
	if tag == language.English {
		return ""
	}
	return "/" + tag.String()
}

func (me *webStaticGoals) slashEmpty(text string) string {
	if text == "" {
		text = "/"
	}
	return text
}

func (me *webStaticGoals) getWebPath(tag language.Tag) string {
	return me.slashEmpty(me.getLanguagePath(tag))
}

func (me *webStaticGoals) getApiPath() string {
	return "http://localhost:8080/hinst-website/pages"
}

func (me *webStaticGoals) getPathQuery(tag language.Tag) string {
	return "&webPath=" + url.QueryEscape(me.getWebPath(tag)) +
		"&apiPath=" + url.QueryEscape(me.getApiPath()) +
		"&staticPath=" + url.QueryEscape("/")
}
