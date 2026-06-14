package server

import (
	"context"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/db_objects"
)

func (me *database) getUrlPings() (results []db_objects.UrlPingRecord) {
	var rows = gophers.AssertResultError(me.pool.Query(context.Background(), "SELECT * FROM urlPings"))
	defer rows.Close()
	for rows.Next() {
		var record db_objects.UrlPingRecord
		record.Scan(rows)
		results = append(results, record)
	}
	return
}

func (me *database) getUrlPing(url string) *db_objects.UrlPingRecord {
	var rows = gophers.AssertResultError(me.pool.Query(context.Background(), "SELECT * FROM urlPings WHERE url = $1", url))
	defer rows.Close()
	if rows.Next() {
		var record db_objects.UrlPingRecord
		record.Scan(rows)
		return &record
	}
	return nil
}

func (me *database) insertUrlPing(url string) {
	gophers.AssertResultError(me.pool.Exec(context.Background(), "INSERT INTO urlPings (url, googlePingedAt) VALUES ($1, NULL)", url))
}

func (me *database) updateUrlPingGoogle(url string, dateTime time.Time) bool {
	var unixSeconds = dateTime.UTC().Unix()
	var result = gophers.AssertResultError(me.pool.Exec(context.Background(), "UPDATE urlPings SET googlePingedAt = $1 WHERE url = $2", unixSeconds, url))
	var rowCount = result.RowsAffected()
	return rowCount > 0
}

func (me *database) updateUrlPingGoogleManually(url string, dateTime time.Time) bool {
	var unixSeconds = dateTime.UTC().Unix()
	var result = gophers.AssertResultError(me.pool.Exec(context.Background(), "UPDATE urlPings SET googlePingedManuallyAt = $1 WHERE url = $2", unixSeconds, url))
	var rowCount = result.RowsAffected()
	return rowCount > 0
}
