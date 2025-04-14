package main

type goalPostObject struct {
	GoalId int64 `json:"goalId"`
	// HTML
	Text string `json:"text"`
	// Unix epoch seconds
	DateTime            int64  `json:"dateTime"`
	IsAutoTranslated    bool   `json:"isAutoTranslated"`
	LanguageName        string `json:"languageName"`
	LanguageNamePending string `json:"languageNamePending"`
	IsPublic            bool   `json:"isPublic"`
}
