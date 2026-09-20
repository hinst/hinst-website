package rest_objects

import (
	"github.com/hinst/hinst-website/server/db_objects"
	"golang.org/x/text/language"
)

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
		me.Title = row.GetTranslatedTitle(language.English)
	}
}

func (me GoalPostHeader) ReadMany(rows []*db_objects.GoalPostRow, languageTag language.Tag) (outputs []*GoalPostHeader) {
	outputs = make([]*GoalPostHeader, 0, len(rows))
	for _, row := range rows {
		var output = new(GoalPostHeader)
		output.Read(row, languageTag)
		outputs = append(outputs, output)
	}
	return
}
