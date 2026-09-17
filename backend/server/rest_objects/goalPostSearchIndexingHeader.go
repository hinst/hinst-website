package rest_objects

import (
	"github.com/hinst/hinst-website/server/db_objects"
	"golang.org/x/text/language"
)

type GoalPostSearchIndexingHeader struct {
	GoalPostHeader
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
