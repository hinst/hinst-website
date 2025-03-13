package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hinst/hinst-website/file_mode"
)

type translator struct {
	savedGoalsPath          string
	translatedDirectoryName string
}

func (me *translator) init() *translator {
	me.translatedDirectoryName = "translated"
	return me
}

func (me *translator) migrate() {
	var goalFiles = assertResultError(os.ReadDir(me.savedGoalsPath))
	for _, goalDirectory := range goalFiles {
		if !goalDirectory.IsDir() {
			continue
		}
		var translatedFilesDir = filepath.Join(
			me.savedGoalsPath, goalDirectory.Name(), me.translatedDirectoryName,
		)
		var files = assertResultError(os.ReadDir(translatedFilesDir))
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			var filePath = filepath.Join(translatedFilesDir, file.Name())
			var newFilePath = getFilePathWithoutExtension(filePath) + ".html"
			log.Println("Migrating " + filePath + " -> " + newFilePath)
			os.Rename(filePath, newFilePath)
		}
	}
}

func (me *translator) run() {
	var goalFiles = assertResultError(os.ReadDir(me.savedGoalsPath))
	for _, goalFile := range goalFiles {
		if !goalFile.IsDir() {
			continue
		}
		var goalFilePath = filepath.Join(me.savedGoalsPath, goalFile.Name())
		me.translateGoal(goalFilePath)
	}
}

func (me *translator) translateGoal(directoryPath string) {
	assertError(os.MkdirAll(
		filepath.Join(directoryPath, "translated"),
		file_mode.OS_USER_RWX,
	))
	var files = assertResultError(os.ReadDir(directoryPath))
	for _, file := range files {
		if !file.IsDir() && GoalFileNameMatcher.MatchString(file.Name()) {
			var filePath = filepath.Join(directoryPath, file.Name())
			me.translateFile(filePath)
		}
	}
}

func (me *translator) getTranslatedFilePath() {
}

func (me *translator) translateFile(filePath string) {
	var article = readJsonFile(filePath, &smartPost{})
	var englishMessage = assertResultError(me.translateText(article.Msg))
	var targetFilePath = filepath.Join(
		filepath.Dir(filePath), "translated",
		getFileNameWithoutExtension(filePath))
	writeTextFile(targetFilePath+".ru.txt", article.Msg)
	writeTextFile(targetFilePath+".en.txt", englishMessage)
}

func (me *translator) translateText(text string) (string, error) {
	var request = encodeJson(lmStudioRequest{
		Model: "aya-expanse-8B",
		Messages: []lmStudioMessage{
			{Role: lm_studio_role_system, Content: prompt_russian_to_english},
			{Role: lm_studio_role_user, Content: text},
		},
		Stream: false,
	})
	var response = assertResultError(http.Post(
		"http://localhost:1235/v1/chat/completions", "application/json", bytes.NewBuffer(request),
	))
	defer func() {
		assertError(response.Body.Close())
	}()
	if response.StatusCode != http.StatusOK {
		return "", errors.New("Cannot translate text. Status: " + response.Status)
	}
	var responseText = assertResultError(io.ReadAll(response.Body))
	var responseObject = decodeJson(responseText, new(lmStudioResponse))
	return responseObject.Choices[0].Message.Content, nil
}
