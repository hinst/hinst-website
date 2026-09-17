package rest_objects

import "github.com/hinst/hinst-website/server/db_objects"

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

func (me *GoalPostObject) Read(goalPostRow *db_objects.GoalPostRow) {
}
