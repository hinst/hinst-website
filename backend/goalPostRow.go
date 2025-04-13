package main

import (
	"database/sql"
	"time"
)

type goalPostRow struct {
	GoalId   int64
	DateTime time.Time
	IsPublic bool
	Text     string
}

func (me *goalPostRow) scan(rows *sql.Rows) {
	var dateTimeMilliseconds int64
	assertError(rows.Scan(&me.GoalId, &dateTimeMilliseconds, &me.IsPublic))
	me.DateTime = time.Unix(dateTimeMilliseconds, 0)
}

func (me *goalPostRow) String() string {
	return "{goalId:" + getStringFromInt64(me.GoalId) +
		", dateTime:" + me.DateTime.String() +
		", isPublic:" + getStringFromBool(me.IsPublic) + "}"
}
