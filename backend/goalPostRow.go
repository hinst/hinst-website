package main

import (
	"database/sql"
	"strconv"
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

func (me *goalPostRow) String() string {
	return "{goalId:" + strconv.Itoa(me.goalId) +
		", dateTime:" + me.dateTime.String() +
		", isPublic:" + getStringFromBool(me.isPublic) + "}"
}
