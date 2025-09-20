package server

import (
	"io"
	"log"
	"net/http"
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
	var path = me.getLanguagePath(lang)
	assertError(os.MkdirAll(path, os.ModePerm))
	log.Printf("path: %v", path)
	var homePage = assertResultError(http.Get(me.url + "/pages?lang=" + lang.String()))
	defer homePage.Body.Close()
	var homePageText = string(assertResultError(io.ReadAll(homePage.Body)))
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
