package main

import (
	"database/sql"
	"strconv"
	"time"
)

type goalPostRow struct {
	GoalId   int
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
	return "{goalId:" + strconv.Itoa(me.GoalId) +
		", dateTime:" + me.DateTime.String() +
		", isPublic:" + getStringFromBool(me.IsPublic) + "}"
}
