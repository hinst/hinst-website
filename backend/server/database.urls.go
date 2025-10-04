package server

func (me *database) getPendingUrlPings() (results []urlPingRecord) {
	var db = me.open()
	defer me.close(db)
	var rows = assertResultError(db.Query("SELECT url, service, doneAt FROM urlPings WHERE doneAt IS NULL"))
	defer ioClose(rows)
	for rows.Next() {
		var record urlPingRecord
		assertError(rows.Scan(&record.url, &record.service, &record.doneAt))
		results = append(results, record)
	}
	return
}
