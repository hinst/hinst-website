package server

import "database/sql"

type urlPingRecord struct {
	Url                    string `json:"url"`
	GooglePingedAt         *int64 `json:"googlePingedAt"`         // Unix seconds UTC
	GooglePingedManuallyAt *int64 `json:"googlePingedManuallyAt"` // Unix seconds UTC
}

func (me *urlPingRecord) scan(rows *sql.Rows) {
	assertError(rows.Scan(
		&me.Url,
		&me.GooglePingedAt,
		&me.GooglePingedManuallyAt,
	))
}
