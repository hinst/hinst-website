package server

import (
	"os"
	"time"

	"github.com/hinst/hinst-website/server/file_mode"
	"golang.org/x/text/language"
)

type webStaticGoals struct {
	folder string
	url    string
	db     *database
}

func (me *webStaticGoals) init(url string, db *database, folder string) {
	me.url = url
	me.db = db
	me.folder = folder
}

func (me *webStaticGoals) run() {
	assertError(os.RemoveAll(me.folder))
	assertError(os.MkdirAll(me.folder, file_mode.OS_USER_RWX))
	assertError(os.CopyFS(me.folder+"/static", os.DirFS("pages/static")))
	for _, lang := range supportedLanguages {
		me.generate(lang)
	}
}

func (me *webStaticGoals) generate(lang language.Tag) {
	var path = me.folder + me.getLanguagePath(lang)
	assertError(os.MkdirAll(path, file_mode.OS_USER_RWX))
	var homeUrl = buildUrl(me.url+"/pages", me.getPathQuery(lang))
	var homePageText = readTextFromUrl(homeUrl)
	writeTextFile(path+"/index.html", homePageText)

	var goals = me.db.getGoals()
	var goalsPath = path + "/personal-goals"
	assertError(os.MkdirAll(goalsPath, file_mode.OS_USER_RWX))
	for _, goal := range goals {
		me.generateGoal(lang, goalsPath, goal)
	}
}

func (me *webStaticGoals) generateGoal(lang language.Tag, goalsPath string, goal goalRecord) {
	var goalId = goal.Id
	var url = buildUrl(me.url+pagesWebPath+"/personal-goals/"+getStringFromInt64(goalId), me.getPathQuery(lang))
	var goalPageText = readTextFromUrl(url)
	writeTextFile(goalsPath+"/"+getStringFromInt64(goalId)+".html", goalPageText)

	var path = goalsPath + "/" + getStringFromInt64(goalId)
	assertError(os.MkdirAll(path, file_mode.OS_USER_RWX))
	var posts = me.db.getGoalPosts(goalId, false, lang)
	for _, post := range posts {
		me.generateGoalPost(lang, goalsPath, goalId, post.DateTime)
	}
}

func (me *webStaticGoals) generateGoalPost(lang language.Tag, goalsPath string, goalId int64, postDateTime int64) {
	var url = buildUrl(me.url+pagesWebPath+"/personal-goals/"+getStringFromInt64(goalId)+"/"+
		getStringFromInt64(postDateTime), me.getPathQuery(lang))
	var postPageText = readTextFromUrl(url)
	var path = goalsPath + "/" + getStringFromInt64(goalId)
	writeTextFile(path+"/"+getStringFromInt64(postDateTime)+".html", postPageText)

	var imageCount = me.db.getGoalPostImageCount(goalId, time.Unix(postDateTime, 0))
	for imageIndex := range imageCount {
		me.generateGoalPostImage(goalId, postDateTime, imageIndex)
	}
}

func (me *webStaticGoals) generateGoalPostImage(goalId int64, postDateTime int64, imageIndex int) {
	var url = buildUrl(me.url+pagesWebPath+"/personal-goals/image/"+getStringFromInt64(goalId)+"/"+
		getStringFromInt64(postDateTime)+"/"+getStringFromInt(imageIndex), nil)
	var image = readBytesFromUrl(url)
	var path = me.folder + "/personal-goals/image/" + getStringFromInt64(goalId) + "/" +
		getStringFromInt64(postDateTime)
	assertError(os.MkdirAll(path, file_mode.OS_USER_RWX))
	path += "/" + getStringFromInt(imageIndex) + ".jpg"
	if checkFileExists(path) {
		return // already saved
	}
	writeBytesFile(path, image)
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

func (me *webStaticGoals) getPathQuery(tag language.Tag) map[string]string {
	return map[string]string{
		"webPath":       me.getWebPath(tag),
		"staticPath":    "/",
		"jpegExtension": ".jpg",
		"htmlExtension": ".html",
		"lang":          tag.String(),
	}
}
