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
	"golang.org/x/text/language"
)

type translator struct {
	savedGoalsPath string
	apiUrl         string
}

const translatedGoalDirectoryName = "translated"

func (me *translator) init() *translator {
	me.apiUrl = "http://localhost:1235/v1/chat/completions"
	return me
}

func (me *translator) migrate() {
	var goalFiles = assertResultError(os.ReadDir(me.savedGoalsPath))
	for _, goalDirectory := range goalFiles {
		if !goalDirectory.IsDir() {
			continue
		}
		var translatedFilesDir = filepath.Join(
			me.savedGoalsPath, goalDirectory.Name(), translatedGoalDirectoryName,
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
		filepath.Join(directoryPath, translatedGoalDirectoryName),
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

func (me *translator) getTranslatedFilePath(smartPostFilePath string, tag language.Tag) string {
	var targetFilePath = filepath.Join(
		filepath.Dir(smartPostFilePath),
		translatedGoalDirectoryName,
		filepath.Base(smartPostFilePath),
	)
	switch tag {
	case language.Russian:
		return targetFilePath + ".ru.html"
	case language.English:
		return targetFilePath + ".en.html"
	default:
		panic(errors.New("Unknown language tag: " + tag.String()))
	}
}

func (me *translator) translateFile(smartPostFilePath string) {
	var article = readJsonFile(smartPostFilePath, &smartPost{})
	var englishMessage = assertResultError(me.translateText(article.Msg))
	writeTextFile(me.getTranslatedFilePath(smartPostFilePath, language.Russian), article.Msg)
	writeTextFile(me.getTranslatedFilePath(smartPostFilePath, language.English), englishMessage)
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
	var response = assertResultError(http.Post(me.apiUrl, contentTypeJson, bytes.NewBuffer(request)))
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
