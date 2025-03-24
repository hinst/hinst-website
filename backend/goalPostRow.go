package main

import (
	"database/sql"
	"time"
)

type goalPostRow struct {
	goalId   int
	dateTime time.Time
	isPublic bool
}

func (me *goalPostRow) scan(rows *sql.Rows) {
	var dateTimeMilliseconds int64
	assertError(rows.Scan(&me.goalId, &dateTimeMilliseconds, &me.isPublic))
	me.dateTime = time.Unix(dateTimeMilliseconds, 0)
}
