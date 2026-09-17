package rest_objects

import (
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"golang.org/x/text/language"
)

type GoalPostObject struct {
	GoalId int64 `json:"goalId"`
	// Unix epoch seconds
	DateTime int64 `json:"dateTime"`
	// HTML
	Text                  string `json:"text"`
	IsAutoTranslated      bool   `json:"isAutoTranslated"`
	IsTranslationPending  bool   `json:"isTranslationPending"`
	LanguageName          string `json:"languageName"`
	LanguageTag           string `json:"languageTag"`
	IsPublic              bool   `json:"isPublic"`
	SearchIndexingEnabled bool   `json:"searchIndexingEnabled,omitempty"`
	ImageCount            int    `json:"imageCount"`
}

func (me *GoalPostObject) Read(goalPostRow *db_objects.GoalPostRow, languageTag language.Tag, isAdminMode bool) {
	me.GoalId = goalPostRow.GoalId
	me.DateTime = goalPostRow.GetDateTime().UTC().Unix()
	me.Text = goalPostRow.GetTranslatedText(base.SupportedLanguages[0])
	me.LanguageTag = languageTag.String()
	me.LanguageName = base.GetLanguageName(languageTag)
	if languageTag != base.SupportedLanguages[0] {
		var translatedText = goalPostRow.GetTranslatedText(languageTag)
		if translatedText != "" {
			me.IsAutoTranslated = true
			me.Text = translatedText
		} else {
			me.IsTranslationPending = true
		}
	}
	me.IsPublic = goalPostRow.IsPublic
	if isAdminMode {
		me.SearchIndexingEnabled = goalPostRow.SearchIndexingEnabled
	}
}
