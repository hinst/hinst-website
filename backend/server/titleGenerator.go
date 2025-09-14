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

type titleGenerator struct {
	apiUrl string
	db     *database
}

var titleGeneratorPreset = titleGenerator{
	apiUrl: lm_studio_default_url,
}

func (me *titleGenerator) run() {
	var totalCount int64
	var updatedCount int64
	me.db.forEachGoalPost(func(row *goalPostRow) bool {
		totalCount++
		var isUpdated = false
		if row.title == nil {
			var title = me.summarizeText(row.text)
			me.db.setGoalPostTitle(row.goalId, row.dateTime, language.Russian, title)
			isUpdated = true
		}
		if row.titleEnglish == nil && row.textEnglish != nil {
			var title = me.summarizeText(*row.textEnglish)
			me.db.setGoalPostTitle(row.goalId, row.dateTime, language.English, title)
			isUpdated = true
		}
		if row.titleGerman == nil && row.textGerman != nil {
			var title = me.summarizeText(*row.textGerman)
			me.db.setGoalPostTitle(row.goalId, row.dateTime, language.German, title)
			isUpdated = true
		}
		if isUpdated {
			updatedCount++
		}
		return true
	})
	var trimmedCount int64
	me.db.forEachGoalPost(func(row *goalPostRow) bool {
		for _, lang := range supportedLanguages {
			var title = row.getTranslatedTitle(lang)
			if title == "" {
				return true
			}
			var trimmedTitle = me.trim(title)
			if title != trimmedTitle {
				trimmedCount++
				me.db.setGoalPostTitle(row.goalId, row.dateTime, lang, trimmedTitle)
			}
		}
		return true
	})
	log.Println("Title generation completed. Total posts:", totalCount,
		"; updated posts:", updatedCount,
		"; trimmed titles:", trimmedCount)
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
	resultText = me.trim(resultText)
	return resultText
}

func (me *titleGenerator) trim(text string) string {
	for {
		var trimmedText = me.trimOnce(text)
		if trimmedText == text {
			return text
		}
		text = trimmedText
	}
}

func (me *titleGenerator) trimOnce(text string) string {
	if strings.HasPrefix(text, "\"") && strings.HasSuffix(text, "\"") {
		text = text[1 : len(text)-1]
	}
	if strings.HasPrefix(text, "*") && strings.HasSuffix(text, "*") {
		text = text[1 : len(text)-1]
	}
	text = strings.TrimPrefix(text, "# ")
	if strings.Contains(text, "\n") {
		var lines = strings.SplitN(text, "\n", 2)
		text = lines[0]
	}
	return text
}
