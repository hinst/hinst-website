package server

import "database/sql"

type urlPingRecord struct {
	url            string
	googlePingedAt *int64 // Unix seconds UTC
}

func (me *urlPingRecord) scan(rows *sql.Rows) {
	assertError(rows.Scan(&me.url, &me.googlePingedAt))
}
