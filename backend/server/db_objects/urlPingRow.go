package db_objects

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/hinst/go-gophers"
	"github.com/jackc/pgx/v5"
)

type UrlPingRow struct {
	Url                    string `json:"url"`
	GooglePingedAt         *int64 `json:"googlePingedAt"`         // Unix seconds UTC
	GooglePingedManuallyAt *int64 `json:"googlePingedManuallyAt"` // Unix seconds UTC
}

var _ = registerDbObject(func() DbObject { return new(UrlPingRow) })

func (UrlPingRow) GetAllColumns() []string {
	return gophers.GetFieldNames[UrlPingRow]()
}

func (UrlPingRow) GetTableName() string {
	return "urlPings"
}

func (me UrlPingRow) SaveToDirectory(directory string) {
	var hash = sha256.Sum256([]byte(me.Url))
	var urlHash = hex.EncodeToString(hash[:])
	var filePath = directory + "/" + urlHash + ".json"
	gophers.WriteJsonFile(filePath, me)
}

func (me *UrlPingRow) Scan(rows pgx.Rows) {
	gophers.AssertError(rows.Scan(
		&me.Url,
		&me.GooglePingedAt,
		&me.GooglePingedManuallyAt,
	))
}
