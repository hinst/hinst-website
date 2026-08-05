package db_objects

import (
	"crypto/sha256"
	"encoding/hex"

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

func (me UrlPingRecord) SaveToDirectory(directory string) {
	var hash = sha256.Sum256([]byte(me.Url))
	var urlHash = hex.EncodeToString(hash[:])
	var filePath = directory + "/" + urlHash + ".json"
	gophers.WriteJsonFile(filePath, me)
}

func (me *UrlPingRecord) Scan(rows pgx.Rows) {
	gophers.AssertError(rows.Scan(
		&me.Url,
		&me.GooglePingedAt,
		&me.GooglePingedManuallyAt,
	))
}
