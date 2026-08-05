package db_objects

import (
	"github.com/hinst/go-gophers"
	"github.com/jackc/pgx/v5"
	"golang.org/x/text/language"
)

type GoalRow struct {
	Id               int64  `json:"id"`
	Title            string `json:"title"`
	TitleEnglish     string `json:"titleEnglish"`
	TitleGerman      string `json:"titleGerman"`
	ImageData        []byte
	ImageContentType string
}

var _ = registerDbObject(GoalRow{})

func (GoalRow) GetTableName() string {
	return "goals"
}

func (me GoalRow) SaveToDirectory(directory string) {
	me.ImageData = nil
	gophers.WriteJsonFile(directory+"/"+gophers.GetStringFromInt64(me.Id), me)
}

func (me *GoalRow) Scan(rows pgx.Rows) {
	gophers.AssertError(rows.Scan(&me.Id, &me.Title, &me.TitleEnglish, &me.TitleGerman,
		&me.ImageData, &me.ImageContentType))
}

func (me GoalRow) GetTranslatedTitle(languageTag language.Tag) string {
	switch languageTag {
	case language.English:
		return me.TitleEnglish
	case language.German:
		return me.TitleGerman
	default:
		return me.Title
	}
}
