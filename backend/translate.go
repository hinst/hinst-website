package main

import (
	"context"
	"os"

	"cloud.google.com/go/translate"
	_ "cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type Translator struct {
	_context        context.Context
	translateClient *translate.Client
}

func (me *Translator) run() {
	var savedGoalsPath = "./saved-goals"
	var goalFiles = assertResultError(os.ReadDir(savedGoalsPath))
	for _, goalFile := range goalFiles {
		if goalFile.IsDir() {
			var goalFilePath = savedGoalsPath + "/" + goalFile.Name()
			me.translateGoal(goalFilePath)
		}
		break
	}
}

func (me *Translator) translateGoal(directoryPath string) {
	var files = assertResultError(os.ReadDir(directoryPath))
	for _, file := range files {
		if !file.IsDir() && GoalFileNameMatcher.MatchString(file.Name()) {
			var filePath = directoryPath + "/" + file.Name()
			me.translateFile(filePath)
		}
		break
	}
}

func (me *Translator) getContext() context.Context {
	if nil == me._context {
		me._context = context.Background()
	}
	return me._context
}

func (me *Translator) getClient() *translate.Client {
	if nil == me.translateClient {
		me.translateClient = assertResultError(translate.NewClient(me.getContext()))
	}
	return me.translateClient
}

func (me *Translator) translateFile(filePath string) {
	var article = readJsonFile(filePath, &smartPost{})
	var translation = assertResultError(me.getClient().Translate(
		me.getContext(),
		[]string{article.Msg},
		language.English,
		&translate.Options{Format: translate.HTML},
	))
	article.Msg = translation[0].Text
	println(article.Msg)
}
