package rest_objects

import (
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"golang.org/x/text/language"
)

type GoalObject struct {
	Id           int64  `json:"id"`
	Title        string `json:"title"`
	TitleEnglish string `json:"titleEnglish"`
	TitleGerman  string `json:"titleGerman"`
}

func (me GoalObject) Read(goalRow db_objects.GoalRow) GoalObject {
	me.Id = goalRow.Id
	me.Title = goalRow.GetTranslatedTitle(base.SupportedLanguages[0])
	me.TitleEnglish = goalRow.GetTranslatedTitle(language.English)
	me.TitleGerman = goalRow.GetTranslatedTitle(language.German)
	return me
}
