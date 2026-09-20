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

const prompt_generate_title = "Generate engaging title for the provided blog post. " +
	"The title should be one sentence in {LANGUAGE} language, plain text. " +
	"Output only the title sentence itself. Commentary stays in the 'thinking' section."

type titleGenerator struct {
	apiUrl string
	db     *database
}

var titleGeneratorPreset = titleGenerator{
	apiUrl: ollama_default_url,
}

func (me *titleGenerator) run() {
	var totalCount int64
	var updatedCount int64
	me.db.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		totalCount++
		var isUpdated = false
		for _, languageTag := range base.SupportedLanguages {
			var translatedText = row.GetTranslatedText(languageTag)
			if row.GetTranslatedTitle(languageTag) == "" && translatedText != "" {
				var title = me.summarizeText(translatedText, languageTag)
				me.db.setGoalPostTitle(row.GoalId, row.GetDateTime(), languageTag, title)
				isUpdated = true
			}
		}
		if isUpdated {
			updatedCount++
		}
		return true
	}, 0)
	log.Printf("Generated title for %v of %v posts\n", updatedCount, totalCount)
}

func (me *titleGenerator) summarizeText(text string, theLanguage language.Tag) string {
	var prompt = strings.ReplaceAll(prompt_generate_title,
		"{LANGUAGE}", base.GetLanguageName(theLanguage))
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
		return errors.New("Cannot summarize text. Status: " + response.Status)
	})
	var responseText = gophers.AssertResultError(io.ReadAll(response.Body))
	var responseObject = gophers.DecodeJson(responseText, new(openAiResponse))
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
