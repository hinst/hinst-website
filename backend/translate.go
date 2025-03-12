package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hinst/hinst-website/file_mode"
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
			var filePath = directoryPath + "/" + file.Name()
			me.translateFile(filePath)
		}
	}
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
