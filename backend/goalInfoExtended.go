package main

type goalHeaderExtended struct {
	goalHeader
	LastPostDate string `json:"lastPostDate"`
}

type smartPostExtended struct {
	smartPost
	IsAutoTranslated    bool   `json:"isAutoTranslated"`
	LanguageTag         string `json:"languageTag"`
	LanguageName        string `json:"languageName"`
	LanguageNamePending string `json:"languageNamePending"`
	IsPublic            bool   `json:"isPublic"`
}

type smartPostHeaderExtended struct {
	smartPostHeader
	IsPublic bool `json:"isPublic"`
}

func (smartPostHeaderExtended) getDatesSeconds(array []*smartPostHeaderExtended) (dates []int64) {
	for _, item := range array {
		var date = assertResultError(parseSmartProgressDate(item.Date))
		dates = append(dates, date.UTC().Unix())
	}
	return
}
