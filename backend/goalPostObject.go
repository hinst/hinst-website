package main

type goalPostObject struct {
	GoalId int64 `json:"goalId"`
	// Unix epoch seconds
	DateTime int64 `json:"dateTime"`
	// HTML
	Text                 string `json:"text"`
	IsAutoTranslated     bool   `json:"isAutoTranslated"`
	IsTranslationPending bool   `json:"isTranslationPending"`
	LanguageName         string `json:"languageName"`
	LanguageTag          string `json:"languageTag"`
	IsPublic             bool   `json:"isPublic"`
}
