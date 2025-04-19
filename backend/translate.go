package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/text/language"
)

type translator struct {
	apiUrl         string
	savedGoalsPath string
	db             *database
}

var translatorPresets = translator{
	apiUrl: "http://localhost:1235/v1/chat/completions",
}

func (me *translator) run() {
	var totalCount = 0
	var translatedCount = 0
	me.db.forEachGoalPost(func(row *goalPostRow) {
		var isDone = false
		if row.textEnglish == nil {
			me.translate(row, language.English)
			isDone = true
		}
		if row.textGerman == nil {
			me.translate(row, language.German)
			isDone = true
		}
		totalCount++
		if isDone {
			translatedCount++
		}
	})
	log.Printf("Translated goal posts: %v of %v", translatedCount, totalCount)
}

func (me *translator) translate(row *goalPostRow, tag language.Tag) {
	var text = assertResultError(me.translateText(row.text, tag))
	me.db.setGoalPostText(row.goalId, row.dateTime, tag, text)
}

func (me *translator) translateText(text string, tag language.Tag) (string, error) {
	var prompt = prompt_Russian_to_something
	prompt = strings.Replace(prompt, "{something}", getLanguageName(tag), -1)
	var request = encodeJson(lmStudioRequest{
		Model: "aya-expanse-8B",
		Messages: []lmStudioMessage{
			{Role: lm_studio_role_system, Content: prompt},
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
