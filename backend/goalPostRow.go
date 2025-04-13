package main

import (
	"database/sql"
	"time"
)

type goalPostRow struct {
	goalId   int64
	dateTime time.Time
	isPublic bool
	text     string

	textEnglish *string
	textGerman  *string
	typeString  string
}

func (me *goalPostRow) scan(rows *sql.Rows) {
	var dateTimeMilliseconds int64
	assertError(rows.Scan(&me.goalId, &dateTimeMilliseconds, &me.isPublic,
		&me.text, &me.textEnglish, &me.textGerman, &me.typeString))
	me.dateTime = time.Unix(dateTimeMilliseconds, 0)
}

func (me *goalPostRow) String() string {
	return "{goalId:" + getStringFromInt64(me.goalId) +
		", dateTime:" + me.dateTime.String() +
		", isPublic:" + getStringFromBool(me.isPublic) + "}"
}
