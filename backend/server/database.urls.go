package server

func (me *database) getPendingUrlPings() (results []urlPingRecord) {
	var db = me.open()
	defer me.close(db)
	var rows = assertResultError(db.Query("SELECT url, googlePingedAt FROM urlPings WHERE googlePingedAt IS NULL"))
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
	var rows = assertResultError(db.Query("SELECT url, googlePingedAt FROM urlPings WHERE googlePingedAt IS NULL"))
	defer ioClose(rows)
	if rows.Next() {
		var record urlPingRecord
		record.scan(rows)
		return &record
	}
	return nil
}
