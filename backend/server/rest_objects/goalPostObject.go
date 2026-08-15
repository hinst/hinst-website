package rest_objects

import (
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

type GoalPostHeader struct {
	GoalId int64 `json:"goalId"`
	// Unix epoch time seconds
	DateTime int64 `json:"dateTime"`
	IsPublic bool  `json:"isPublic"`
	// "post" or "comment"
	Type  string `json:"type"`
	Title string `json:"title"`
}

func (me *GoalPostHeader) Read(row *db_objects.GoalPostRow, languageTag language.Tag) {
	me.GoalId = row.GoalId
	me.DateTime = row.DateTime
	me.IsPublic = row.IsPublic
	me.Type = row.Type
	me.Title = row.GetTranslatedTitle(languageTag)
	if me.Title == "" {
		me.Title = row.TitleEnglish
	}
}

type GoalPostSearchIndexingHeader struct {
	// `tstype:",extends"` makes tygo generate a TypeScript interface that
	// extends `GoalPostHeader` instead of a nested field.
	GoalPostHeader `tstype:",extends"`
	GooglePingedAt                      int64  `json:"googlePingedAt"`
	PublicUrl                           string `json:"publicUrl"`
	GoogleSearchIndexingStatus          string `json:"googleSearchIndexingStatus"`
	GoogleSearchIndexingStatusCheckedAt int64  `json:"googleSearchIndexingStatusCheckedAt"`
}

// Does not set PublicUrl
func (me *GoalPostSearchIndexingHeader) Read(row *db_objects.GoalPostRow, languageTag language.Tag) {
	me.GoalPostHeader.Read(row, languageTag)
	me.GooglePingedAt = row.GooglePingedAt
	me.GoogleSearchIndexingStatus = row.GoogleSearchIndexingStatus
	me.GoogleSearchIndexingStatusCheckedAt = row.GoogleSearchIndexingStatusCheckedAt
}
