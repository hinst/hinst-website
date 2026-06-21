package server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/db_objects"
	"golang.org/x/text/language"
)

const prompt_generate_title = "You are a professional blog post editor. " +
	"Your task is to create an engaging title for the following text. " +
	"The title should be one sentence in plain text. " +
	"Do not provide any explanations or commentary - just the title in the same language as the original text."

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
	me.db.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		totalCount++
		var isUpdated = false
		if row.Title == nil {
			var title = me.summarizeText(row.Text)
			me.db.setGoalPostTitle(row.GoalId, row.GetDateTime(), language.Russian, title)
			isUpdated = true
		}
		if row.TitleEnglish == nil && row.TextEnglish != nil {
			var title = me.summarizeText(*row.TextEnglish)
			me.db.setGoalPostTitle(row.GoalId, row.GetDateTime(), language.English, title)
			isUpdated = true
		}
		if row.TitleGerman == nil && row.TextGerman != nil {
			var title = me.summarizeText(*row.TextGerman)
			me.db.setGoalPostTitle(row.GoalId, row.GetDateTime(), language.German, title)
			isUpdated = true
		}
		if isUpdated {
			updatedCount++
		}
		return true
	}, "*", 0)
	log.Printf("Generated title for %v of %v posts\n", updatedCount, totalCount)
}

func (me *titleGenerator) summarizeText(text string) string {
	var request = gophers.EncodeJson(lmStudioRequest{
		Model: lm_studio_multilingual_model_id,
		Messages: []lmStudioMessage{
			{Role: lm_studio_role_system, Content: prompt_generate_title},
			{Role: lm_studio_role_user, Content: text},
		},
		Stream: false,
	})
	var response = gophers.AssertResultError(http.Post(me.apiUrl, gophers.ContentTypeJson, bytes.NewBuffer(request)))
	defer gophers.IoCloseSilently(response.Body)
	gophers.AssertCondition(response.StatusCode == http.StatusOK, func() error {
		return errors.New("Cannot summarize text. Status: " + response.Status)
	})
	var responseText = gophers.AssertResultError(io.ReadAll(response.Body))
	var responseObject = gophers.DecodeJson(responseText, new(lmStudioResponse))
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
