package server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hinst/go-common"
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
	}, "*", 0)
	log.Printf("Generated translated text for %v of %v posts", translatedCount, totalCount)
}

func (me *translator) migrate() {
	var totalCount = 0
	var translatedCount = 0
	me.db.forEachGoalPost(func(row *goalPostRow) bool {
		log.Println("Checking...")
		var isDone = false
		if row.textEnglish != nil {
			var e = validateHtml(*row.textEnglish)
			if e != nil {
				log.Printf("English text is not valid HTML. Regenerating translation for goalId=%v, dateTime=%v. Error: %v",
					row.goalId, row.dateTime, e)
				me.translate(row, language.English)
				isDone = true
			}
		}
		if row.textGerman != nil {
			var e = validateHtml(*row.textGerman)
			if e != nil {
				log.Printf("German text is not valid HTML. Regenerating translation for goalId=%v, dateTime=%v. Error: %v",
					row.goalId, row.dateTime, e)
				me.translate(row, language.German)
				isDone = true
			}
		}
		totalCount++
		if isDone {
			translatedCount++
		}
		return true
	}, "*", 0)
	log.Printf("Regenerated translated text for %v of %v posts", translatedCount, totalCount)
}

func (me *translator) translate(row *goalPostRow, tag language.Tag) {
	var text = ""
	const attemptLimit = 30
	for i := range attemptLimit {
		var text = common.AssertResultError(formatHtml(row.text))
		text = me.translateText(text, tag)
		var e = validateHtml(text)
		if e == nil {
			break
		} else if i == attemptLimit-1 {
			log.Printf("Cannot generate valid HTML text after %v attempts for goalId=%v, dateTime=%v, language=%v. Last error: %v",
				attemptLimit, row.goalId, row.dateTime, tag, e)
		}
	}
	me.db.setGoalPostText(row.goalId, row.getDateTime(), tag, text)
}

func (me *translator) translateText(text string, tag language.Tag) string {
	var prompt = prompt_Russian_to_something
	prompt = strings.ReplaceAll(prompt, "{something}", getLanguageName(tag))
	var request = common.EncodeJson(lmStudioRequest{
		Model: lm_studio_multilingual_model_id,
		Messages: []lmStudioMessage{
			{Role: lm_studio_role_system, Content: prompt},
			{Role: lm_studio_role_user, Content: text},
		},
		Stream: false,
	})
	var client = &http.Client{Timeout: 1 * time.Hour}
	var req = common.AssertResultError(http.NewRequest("POST", me.apiUrl, bytes.NewBuffer(request)))
	req.Header.Set(common.ContentTypeHeader, common.ContentTypeJson)
	var response = common.AssertResultError(client.Do(req))
	defer common.IoCloseSilently(response.Body)
	common.AssertCondition(response.StatusCode == http.StatusOK, func() error {
		return errors.New("Cannot translate text. Status: " + response.Status)
	})
	var responseText = common.AssertResultError(io.ReadAll(response.Body))
	var responseObject = common.DecodeJson(responseText, new(lmStudioResponse))
	return responseObject.Choices[0].Message.Content
}
