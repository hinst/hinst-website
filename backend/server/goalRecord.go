package server

import (
	"github.com/hinst/go-common"
	"github.com/jackc/pgx/v5"
)

type goalRecord struct {
	Id               int64  `json:"id"`
	Title            string `json:"title"`
	ImageData        []byte
	ImageContentType string
}

func (me *goalRecord) scan(rows pgx.Rows) {
	common.AssertError(rows.Scan(&me.Id, &me.Title, &me.ImageData, &me.ImageContentType))
}

type goalPostRecord struct {
	GoalId int64 `json:"goalId"`
	// Unix epoch time seconds
	DateTime int64 `json:"dateTime"`
	IsPublic bool  `json:"isPublic"`
	// "post" or "comment"
	Type  string  `json:"type"`
	Title *string `json:"title"`
}
