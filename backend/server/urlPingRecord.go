package server

import "github.com/jackc/pgx/v5"

type urlPingRecord struct {
	Url                    string `json:"url"`
	GooglePingedAt         *int64 `json:"googlePingedAt"`         // Unix seconds UTC
	GooglePingedManuallyAt *int64 `json:"googlePingedManuallyAt"` // Unix seconds UTC
}

func (me *urlPingRecord) scan(rows pgx.Rows) {
	AssertError(rows.Scan(
		&me.Url,
		&me.GooglePingedAt,
		&me.GooglePingedManuallyAt,
	))
}
