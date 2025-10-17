package server

import "time"

func (me *database) getPendingUrlPings() (results []urlPingRecord) {
	var db = me.open()
	defer me.close(db)
	var rows = assertResultError(db.Query("SELECT * FROM urlPings WHERE googlePingedAt IS NULL"))
	defer ioClose(rows)
	for rows.Next() {
		var record urlPingRecord
		record.scan(rows)
		results = append(results, record)
	}
	return
}

func (me *database) getUrlPing(url string) *urlPingRecord {
	var db = me.open()
	defer me.close(db)
	var rows = assertResultError(db.Query("SELECT * FROM urlPings WHERE url = ?", url))
	defer ioClose(rows)
	if rows.Next() {
		var record urlPingRecord
		record.scan(rows)
		return &record
	}
	return nil
}

func (me *database) insertUrlPing(url string) {
	var db = me.open()
	defer me.close(db)
	var result = assertResultError(db.Exec("INSERT INTO urlPings (url, googlePingedAt) VALUES (?, NULL)", url))
	assertResultError(result.RowsAffected())
}

func (me *database) updateUrlPingGoogle(url string, dateTime time.Time) bool {
	var db = me.open()
	defer me.close(db)
	var unixSeconds = dateTime.UTC().Unix()
	var result = assertResultError(db.Exec("UPDATE urlPings SET googlePingedAt = ? WHERE url = ?", unixSeconds, url))
	var rowCount = assertResultError(result.RowsAffected())
	return rowCount > 0
}

func (me *database) updateUrlPingGoogleManually(url string, dateTime time.Time) bool {
	var db = me.open()
	defer me.close(db)
	var unixSeconds = dateTime.UTC().Unix()
	var result = assertResultError(db.Exec("UPDATE urlPings SET googlePingedManuallyAt = ? WHERE url = ?", unixSeconds, url))
	var rowCount = assertResultError(result.RowsAffected())
	return rowCount > 0
}
