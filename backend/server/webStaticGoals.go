package server

import (
	"log"
	"os"

	"golang.org/x/text/language"
)

type webStaticGoals struct {
	folder string
	url    string
}

func (me *webStaticGoals) init(url string) {
	me.folder = "static"
	me.url = url
}

func (me *webStaticGoals) run() {
	assertError(os.RemoveAll(me.folder))
	assertError(os.MkdirAll(me.folder, os.ModePerm))
	for _, lang := range supportedLanguages {
		me.generate(lang)
	}
}

func (me *webStaticGoals) generate(lang language.Tag) {
	os.CopyFS("static/hinst-website/pages/static", os.DirFS("pages/static"))

	var path = me.getLanguagePath(lang)
	assertError(os.MkdirAll(path, os.ModePerm))
	log.Printf("path: %v", path)
	var homePageText = readTextFromUrl(me.url + "/pages?lang=" + lang.String())
	writeTextFile(path+"/index.html", homePageText)
}

func (me *webStaticGoals) getLanguagePath(tag language.Tag) (path string) {
	path = me.folder
	if tag == language.English {
		return
	}
	path += "/" + tag.String()
	return
}
