package server

import (
	"github.com/hinst/go-common"
	"github.com/jackc/pgx/v5"
	"golang.org/x/text/language"
)

type goalRow struct {
	Id               int64  `json:"id"`
	Title            string `json:"title"`
	TitleEnglish     string `json:"titleEnglish"`
	TitleGerman      string `json:"titleGerman"`
	ImageData        []byte
	ImageContentType string
}

func (me *goalRow) scan(rows pgx.Rows) {
	common.AssertError(rows.Scan(&me.Id, &me.Title, &me.TitleEnglish, &me.TitleGerman,
		&me.ImageData, &me.ImageContentType))
}

func (me goalRow) getTranslatedTitle(languageTag language.Tag) string {
	switch languageTag {
	case language.English:
		return me.TitleEnglish
	case language.German:
		return me.TitleGerman
	default:
		return me.Title
	}
}
