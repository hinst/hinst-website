package db_objects

import (
	"github.com/hinst/go-gophers"
	"github.com/jackc/pgx/v5"
)

type UrlPingRecord struct {
	Url                    string `json:"url"`
	GooglePingedAt         *int64 `json:"googlePingedAt"`         // Unix seconds UTC
	GooglePingedManuallyAt *int64 `json:"googlePingedManuallyAt"` // Unix seconds UTC
}

var _ = registerDbObject(func() DbObject { return new(UrlPingRecord) })

func (UrlPingRecord) GetAllColumns() []string {
	return gophers.GetFieldNames[UrlPingRecord]()
}

func (UrlPingRecord) GetTableName() string {
	return "urlPings"
}

func (UrlPingRecord) SaveToDirectory(directory string) {
	//TODO
}

func (me *UrlPingRecord) Scan(rows pgx.Rows) {
	gophers.AssertError(rows.Scan(
		&me.Url,
		&me.GooglePingedAt,
		&me.GooglePingedManuallyAt,
	))
}
