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
		var isChanged = false
		for languageIndex := 1; languageIndex < len(base.SupportedLanguages); languageIndex++ {
			var languageTag = base.SupportedLanguages[languageIndex]
			if row.GetTranslatedText(languageTag) == "" && row.GetTranslatedText(me.defaultLanguage()) != "" {
				me.translate(row, languageTag)
				isChanged = true
			}
		}
		totalCount++
		if isChanged {
			translatedCount++
		}
		return true
	}, (db_objects.GoalPostRow{}).GetAllFieldSelector(), 0)
	log.Printf("Generated translated text for %v of %v posts", translatedCount, totalCount)
}

func (me *translator) translate(row *db_objects.GoalPostRow, tag language.Tag) {
	var text = me.translateText(row.GetTranslatedText(me.defaultLanguage()), tag)
	me.db.setGoalPostText(row.GoalId, row.GetDateTime(), tag, text)
}

func (me *translator) translateText(text string, tag language.Tag) string {
	var prompt = PROMPT_TRANSLATE_INTO_LANGUAGE
	prompt = strings.ReplaceAll(prompt, "{LANGUAGE}", base.GetLanguageName(tag))
	var requestObject = gophers.EncodeJson(openAiRequest{
		Model: ollama_model_id,
		Messages: []openAiMessage{
			{Role: AI_ROLE_SYSTEM, Content: prompt},
			{Role: AI_ROLE_USER, Content: text},
		},
		Stream: false,
	})
	var requestFactory = func() *http.Request {
		var requestHttp = gophers.AssertResultError(
			http.NewRequest(http.MethodPost, me.apiUrl, bytes.NewBuffer(requestObject)))
		requestHttp.Header.Set(gophers.ContentTypeHeader, gophers.ContentTypeJson)
		return requestHttp
	}
	var response = gophers.AssertResultError(
		gophers.WebRetry{}.Run(&http.Client{Timeout: 1 * time.Hour}, requestFactory))
	defer gophers.IoCloseSilently(response.Body)
	gophers.AssertCondition(response.StatusCode == http.StatusOK, func() error {
		return errors.New("Cannot translate text. Status: " + response.Status)
	})
	var responseText = gophers.AssertResultError(io.ReadAll(response.Body))
	var responseObject = gophers.DecodeJson(responseText, new(openAiResponse))
	return responseObject.Choices[0].Message.Content
}

func (me translator) defaultLanguage() language.Tag {
	return base.SupportedLanguages[0]
}
