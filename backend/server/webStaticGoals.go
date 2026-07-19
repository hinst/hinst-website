package server

import (
	"os"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/go-gophers/file_mode"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"golang.org/x/text/language"
)

type webStaticGoals struct {
	folder   string
	url      string
	db       *database
	renderer *goalRenderer
}

func (me *webStaticGoals) init(url string, db *database, folder string) {
	me.url = url
	me.db = db
	me.folder = folder
	me.renderer = newGoalRenderer(db, "/pages")
}

func (me *webStaticGoals) run() {
	gophers.AssertError(os.MkdirAll(me.folder, file_mode.OS_USER_RWX))
	me.deleteOldFiles()
	gophers.AssertError(os.CopyFS(me.folder+"/static", os.DirFS("pages/static")))
	for _, lang := range base.SupportedLanguages {
		me.generate(lang)
	}
}

func (me *webStaticGoals) deleteOldFiles() {
	for _, file := range gophers.AssertResultError(os.ReadDir(me.folder)) {
		var filePath = me.folder + "/" + file.Name()
		var isPreserved = staticFilesUpdate{}.checkPreservedFile(file.Name())
		if !isPreserved {
			gophers.AssertError(os.RemoveAll(filePath))
		}
	}
}

func (me *webStaticGoals) generate(lang language.Tag) {
	var path = me.folder + me.getLanguagePath(lang)
	gophers.AssertError(os.MkdirAll(path, file_mode.OS_USER_RWX))

	var req = WebRequest{
		Language:      lang,
		WebPath:       me.getWebPath(lang),
		StaticPath:    "",
		JpegExtension: ".jpg",
		HtmlExtension: ".html",
	}
	var homePageText = me.renderer.renderHomePage(req)
	gophers.WriteTextFile(path+"/index.html", gophers.AssertResultError(formatHtml(homePageText)))

	var goals = me.db.getGoals()
	var goalsPath = path + "/personal-goals"
	gophers.AssertError(os.MkdirAll(goalsPath, file_mode.OS_USER_RWX))
	for _, goal := range goals {
		me.generateGoal(lang, goalsPath, goal)
	}
}

func (me *webStaticGoals) generateGoal(lang language.Tag, goalsPath string, goal db_objects.GoalRow) {
	var goalId = goal.Id

	var req = WebRequest{
		Language:      lang,
		WebPath:       me.getWebPath(lang),
		StaticPath:    "",
		JpegExtension: ".jpg",
		HtmlExtension: ".html",
	}
	var goalPageText = me.renderer.renderGoalPage(req, goalId)
	gophers.WriteTextFile(
		goalsPath+"/"+gophers.GetStringFromInt64(goalId)+".html",
		gophers.AssertResultError(formatHtml(goalPageText)))

	var path = goalsPath + "/" + gophers.GetStringFromInt64(goalId)
	gophers.AssertError(os.MkdirAll(path, file_mode.OS_USER_RWX))
	var posts = me.db.getGoalPosts(goalId, false, lang)
	for _, post := range posts {
		me.generateGoalPost(lang, goalsPath, goalId, post.DateTime)
	}
}

func (me *webStaticGoals) generateGoalPost(lang language.Tag, goalsPath string, goalId int64, postDateTime int64) {
	var req = WebRequest{
		Language:      lang,
		WebPath:       me.getWebPath(lang),
		StaticPath:    "",
		JpegExtension: ".jpg",
		HtmlExtension: ".html",
	}
	var postPageText = me.renderer.renderGoalPostPage(req, goalId, time.Unix(postDateTime, 0))
	var postFilePath = goalsPath + "/" + gophers.GetStringFromInt64(goalId) + "/"
	gophers.WriteTextFile(
		postFilePath+gophers.GetStringFromInt64(postDateTime)+".html",
		gophers.AssertResultError(formatHtml(postPageText)))

	var imageCount = me.db.getGoalPostImageCount(goalId, time.Unix(postDateTime, 0))
	for imageIndex := range imageCount {
		me.generateGoalPostImage(goalId, postDateTime, imageIndex)
	}
}

func (me *webStaticGoals) generateGoalPostImage(goalId int64, postDateTime int64, imageIndex int) {
	var image = me.db.getGoalPostImage(goalId, time.Unix(postDateTime, 0), imageIndex)

	if image == nil {
		return // should not happen based on imageCount
	}
	var path = me.folder + "/personal-goals/image/" + gophers.GetStringFromInt64(goalId) + "/" +
		gophers.GetStringFromInt64(postDateTime)
	gophers.AssertError(os.MkdirAll(path, file_mode.OS_USER_RWX))
	path += "/" + gophers.GetStringFromInt(imageIndex) + ".jpg"
	if gophers.CheckFileExists(path) {
		return // already saved
	}
	gophers.WriteBytesFile(path, image.File)
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
