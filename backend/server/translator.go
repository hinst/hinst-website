package server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"golang.org/x/text/language"
)

const PROMPT_TRANSLATE_INTO_LANGUAGE = "Translate provided blog post into {LANGUAGE} language. " +
	"Output only the translated text itself. Commentary stays in the 'thinking' section."

type translator struct {
	apiUrl string
	db     *database
}

func (me *translator) run() {
	var totalCount = 0
	var translatedCount = 0
	me.db.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		var isDone = false
		if row.TextEnglish == nil {
			me.translate(row, language.English)
			isDone = true
		}
		if row.TextGerman == nil {
			me.translate(row, language.German)
			isDone = true
		}
		totalCount++
		if isDone {
			translatedCount++
		}
		return true
	}, "*", 0)
	log.Printf("Generated translated text for %v of %v posts", translatedCount, totalCount)
}

func (me *translator) translate(row *db_objects.GoalPostRow, tag language.Tag) {
	var text = me.translateText(row.Text, tag)
	me.db.setGoalPostText(row.GoalId, row.GetDateTime(), tag, &text)
}

func (me *translator) translateText(text string, tag language.Tag) string {
	var prompt = PROMPT_TRANSLATE_INTO_LANGUAGE
	prompt = strings.ReplaceAll(prompt, "{LANGUAGE}", base.GetLanguageName(tag))
	var request = gophers.EncodeJson(openAiRequest{
		Model: ollama_model_id,
		Messages: []openAiMessage{
			{Role: AI_ROLE_SYSTEM, Content: prompt},
			{Role: AI_ROLE_USER, Content: text},
		},
		Stream: false,
	})
	var client = &http.Client{Timeout: 1 * time.Hour}
	var req = gophers.AssertResultError(http.NewRequest("POST", me.apiUrl, bytes.NewBuffer(request)))
	req.Header.Set(gophers.ContentTypeHeader, gophers.ContentTypeJson)
	var response = gophers.AssertResultError(client.Do(req))
	defer gophers.IoCloseSilently(response.Body)
	gophers.AssertCondition(response.StatusCode == http.StatusOK, func() error {
		return errors.New("Cannot translate text. Status: " + response.Status)
	})
	var responseText = gophers.AssertResultError(io.ReadAll(response.Body))
	var responseObject = gophers.DecodeJson(responseText, new(openAiResponse))
	return responseObject.Choices[0].Message.Content
}
