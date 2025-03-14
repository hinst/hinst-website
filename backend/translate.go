package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/hinst/hinst-website/file_mode"
	"golang.org/x/text/language"
)

type translator struct {
	apiUrl                string
	translatedGoalDirName string
	supportedLanguages    []language.Tag

	savedGoalsPath string
}

var translatorPresets = translator{
	apiUrl:                "http://localhost:1235/v1/chat/completions",
	translatedGoalDirName: "translated",
	supportedLanguages:    []language.Tag{language.Russian, language.English, language.German},
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
		filepath.Join(directoryPath, me.translatedGoalDirName),
		file_mode.OS_USER_RWX,
	))
	var files = assertResultError(os.ReadDir(directoryPath))
	for _, file := range files {
		if !file.IsDir() && GoalFileNameMatcher.MatchString(file.Name()) {
			var filePath = filepath.Join(directoryPath, file.Name())
			var translatedFilePath = me.getTranslatedFilePath(filePath, language.English)
			if !checkFileExists(translatedFilePath) {
				me.translateFile(filePath)
			}
		}
	}
}

func (me *translator) getTranslatedFilePath(smartPostFilePath string, tag language.Tag) string {
	if !slices.Contains(me.supportedLanguages, tag) {
		panic(errors.New("Language is not supported: " + tag.String()))
	}
	var targetFilePath = filepath.Join(
		filepath.Dir(smartPostFilePath),
		me.translatedGoalDirName,
		getFileNameWithoutExtension(smartPostFilePath),
	)
	return targetFilePath + "." + tag.String() + ".html"
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
