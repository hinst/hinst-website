package main

type goalHeaderExtended struct {
	goalHeader
	LastPostDate string `json:"lastPostDate"`
}

type smartPostExtended struct {
	smartPost
	FileName            string `json:"fileName"`
	IsAutoTranslated    bool   `json:"isAutoTranslated"`
	LanguageTag         string `json:"languageTag"`
	LanguageName        string `json:"languageName"`
	LanguageNamePending string `json:"languageNamePending"`
}
