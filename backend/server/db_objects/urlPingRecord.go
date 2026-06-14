package db_objects

import (
	"github.com/hinst/go-common"
	"github.com/jackc/pgx/v5"
)

type UrlPingRecord struct {
	Url                    string `json:"url"`
	GooglePingedAt         *int64 `json:"googlePingedAt"`         // Unix seconds UTC
	GooglePingedManuallyAt *int64 `json:"googlePingedManuallyAt"` // Unix seconds UTC
}

func (me *UrlPingRecord) Scan(rows pgx.Rows) {
	common.AssertError(rows.Scan(
		&me.Url,
		&me.GooglePingedAt,
		&me.GooglePingedManuallyAt,
	))
}
