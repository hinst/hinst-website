package server

import (
	"github.com/hinst/go-common"
	"github.com/jackc/pgx/v5"
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

type goalPostHeader struct {
	GoalId int64 `json:"goalId"`
	// Unix epoch time seconds
	DateTime int64 `json:"dateTime"`
	IsPublic bool  `json:"isPublic"`
	// "post" or "comment"
	Type  string  `json:"type"`
	Title *string `json:"title"`
}
