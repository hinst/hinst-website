package rest_objects

import "github.com/hinst/hinst-website/server/db_objects"

type GoalObject struct {
	Id           int64  `json:"id"`
	Title        string `json:"title"`
	TitleEnglish string `json:"titleEnglish"`
	TitleGerman  string `json:"titleGerman"`
}

func (me GoalObject) Read(goalRow db_objects.GoalRow) GoalObject {
	me.Id = goalRow.Id
	me.Title = goalRow.Title
	me.TitleEnglish = goalRow.TitleEnglish
	me.TitleGerman = goalRow.TitleGerman
	return me
}
