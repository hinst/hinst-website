package server

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"golang.org/x/text/language"
)

type titleGenerator struct {
	apiUrl string
	db     *database
}

var titleGeneratorPreset = titleGenerator{
	apiUrl: lm_studio_default_url,
}

func (me *titleGenerator) run() {
	me.db.forEachGoalPost(func(row *goalPostRow) bool {
		if row.title == nil {
			var title = me.summarizeText(row.text)
			me.db.setGoalPostTitle(row.goalId, row.dateTime, language.Russian, title)
		}
		if row.titleEnglish == nil && row.textEnglish != nil {
			var title = me.summarizeText(*row.textEnglish)
			me.db.setGoalPostTitle(row.goalId, row.dateTime, language.English, title)
		}
		if row.titleGerman == nil && row.textGerman != nil {
			var title = me.summarizeText(*row.textGerman)
			me.db.setGoalPostTitle(row.goalId, row.dateTime, language.German, title)
		}
		return true
	})
}

func (me *titleGenerator) summarizeText(text string) string {
	var request = encodeJson(lmStudioRequest{
		Model: lm_studio_multilingual_model_id,
		Messages: []lmStudioMessage{
			{Role: lm_studio_role_system, Content: prompt_generate_title},
			{Role: lm_studio_role_user, Content: text},
		},
		Stream: false,
	})
	var response = assertResultError(http.Post(me.apiUrl, contentTypeJson, bytes.NewBuffer(request)))
	defer response.Body.Close()
	assertCondition(response.StatusCode == http.StatusOK, func() error {
		return errors.New("Cannot summarize text. Status: " + response.Status)
	})
	var responseText = assertResultError(io.ReadAll(response.Body))
	var responseObject = decodeJson(responseText, new(lmStudioResponse))
	var resultText = responseObject.Choices[0].Message.Content
	if len(resultText) > 1 && resultText[0] == '"' && resultText[len(resultText)-1] == '"' {
		resultText = resultText[1 : len(resultText)-1]
	}
	return resultText
}
