package server

import (
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
	for _, lang := range supportedLanguages {
		me.generate(lang)
	}
}

func (me *webStaticGoals) generate(lang language.Tag) {
	os.CopyFS(me.folder+"/pages/static", os.DirFS("pages/static"))

	var path = me.getLanguagePath(lang)
	assertError(os.MkdirAll(path, os.ModePerm))
	var homePageText = readTextFromUrl(me.url + "/pages?lang=" + lang.String())
	writeTextFile(path+"/index.html", homePageText)

	var goals = me.db.getGoals()
	var goalsPath = path + pagesWebPath + "/personal-goals"
	assertError(os.MkdirAll(goalsPath, os.ModePerm))
	for _, goal := range goals {
		me.generateGoal(lang, goalsPath, goal)
	}
}

func (me *webStaticGoals) generateGoal(lang language.Tag, goalsPath string, goal goalRecord) {
	var goalId = goal.Id
	var url = me.url + pagesWebPath + "/personal-goals/" + getStringFromInt64(goalId) + "?lang=" + lang.String()
	var goalPageText = readTextFromUrl(url)
	writeTextFile(goalsPath+"/"+getStringFromInt64(goalId)+".html", goalPageText)
}

func (me *webStaticGoals) getLanguagePath(tag language.Tag) (path string) {
	path = me.folder
	if tag == language.English {
		return
	}
	path += "/" + tag.String()
	return
}
