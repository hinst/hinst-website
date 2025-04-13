package main

import "time"

type goalPostObject struct {
	GoalId int64 `json:"goalId"`
	// HTML
	Text                string    `json:"text"`
	DateTime            time.Time `json:"dateTime"`
	IsAutoTranslated    bool      `json:"isAutoTranslated"`
	LanguageName        string    `json:"languageName"`
	LanguageNamePending string    `json:"languageNamePending"`
}
