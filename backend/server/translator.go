package server

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
	apiUrl string
	db     *database
}

var translatorPreset = translator{
	apiUrl: lm_studio_default_url,
}

func (me *translator) run() {
	var totalCount = 0
	var translatedCount = 0
	me.db.forEachGoalPost(func(row *goalPostRow) bool {
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
		return true
	})
	log.Printf("Translated goal posts: %v of %v", translatedCount, totalCount)
}

func (me *translator) translate(row *goalPostRow, tag language.Tag) {
	var text = me.translateText(row.text, tag)
	me.db.setGoalPostText(row.goalId, row.dateTime, tag, text)
}

func (me *translator) translateText(text string, tag language.Tag) string {
	var prompt = prompt_Russian_to_something
	prompt = strings.ReplaceAll(prompt, "{something}", getLanguageName(tag))
	var request = encodeJson(lmStudioRequest{
		Model: lm_studio_multilingual_model_id,
		Messages: []lmStudioMessage{
			{Role: lm_studio_role_system, Content: prompt},
			{Role: lm_studio_role_user, Content: text},
		},
		Stream: false,
	})
	var response = assertResultError(http.Post(me.apiUrl, contentTypeJson, bytes.NewBuffer(request)))
	defer response.Body.Close()
	assertCondition(response.StatusCode == http.StatusOK, func() error {
		return errors.New("Cannot translate text. Status: " + response.Status)
	})
	var responseText = assertResultError(io.ReadAll(response.Body))
	var responseObject = decodeJson(responseText, new(lmStudioResponse))
	return responseObject.Choices[0].Message.Content
}
