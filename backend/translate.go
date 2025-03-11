package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
)

type translator struct {
}

func (me *translator) run() {
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

func (me *translator) translateGoal(directoryPath string) {
	var files = assertResultError(os.ReadDir(directoryPath))
	for _, file := range files {
		if !file.IsDir() && GoalFileNameMatcher.MatchString(file.Name()) {
			var filePath = directoryPath + "/" + file.Name()
			me.translateFile(filePath)
		}
		break
	}
}

func (me *translator) translateFile(filePath string) {
	var article = readJsonFile(filePath, &smartPost{})
	var request = encodeJson(libreTranslateRequest{
		Query:  article.Msg,
		Source: "ru",
		Target: "en",
		Format: "html",
	})
	var response = assertResultError(http.Post(
		"http://localhost:5000/translate", "application/json", bytes.NewBuffer(request),
	))
	defer func() {
		assertError(response.Body.Close())
	}()
	var responseText = assertResultError(io.ReadAll(response.Body))
	log.Printf("Response text: %s\n", responseText)
	var translatedArticle = decodeJson(responseText, new(libreTranslateResponse))
	log.Printf("Translated article: %s\n", translatedArticle.TranslatedText)
}
