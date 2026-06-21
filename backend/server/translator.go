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

const prompt_Russian_to_something = "You are a professional Russian-to-{something} translator specializing in diary blog posts. " +
	" Your task is to provide accurate, contextually appropriate translations while preserving markdown formatting. " +
	" Do not provide any explanations or commentary - just the direct translation."

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
	var text = ""
	const attemptLimit = 30
	for i := range attemptLimit {
		text = gophers.AssertResultError(formatHtml(row.Text))
		text = me.translateText(text, tag)
		var e = validateHtml(text)
		if e == nil {
			break
		} else if i == attemptLimit-1 {
			log.Printf("Cannot generate valid HTML text after %v attempts for goalId=%v, dateTime=%v, language=%v. Last error: %v",
				attemptLimit, row.GoalId, row.DateTime, tag, e)
		}
	}
	me.db.setGoalPostText(row.GoalId, row.GetDateTime(), tag, &text)
}

func (me *translator) translateText(text string, tag language.Tag) string {
	var prompt = prompt_Russian_to_something
	prompt = strings.ReplaceAll(prompt, "{something}", base.GetLanguageName(tag))
	var request = gophers.EncodeJson(lmStudioRequest{
		Model: lm_studio_multilingual_model_id,
		Messages: []lmStudioMessage{
			{Role: lm_studio_role_system, Content: prompt},
			{Role: lm_studio_role_user, Content: text},
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
	var responseObject = gophers.DecodeJson(responseText, new(lmStudioResponse))
	return responseObject.Choices[0].Message.Content
}
