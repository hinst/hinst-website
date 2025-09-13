package server

import (
	"database/sql"
	"strconv"
	"time"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

// Measured in milliseconds
var dbTimeout = time.Hour / time.Millisecond

//go:embed schema.sql
var dbSchema string

type database struct {
	dataDirectory string
}

func (me *database) init(dataDirectory string) {
	me.dataDirectory = dataDirectory
	var db = me.open()
	defer me.close(db)
	assertError(db.Ping())
	assertResultError(db.Exec(dbSchema))
	me.collectGarbage()
}

func (me *database) open() *sql.DB {
	return me.openFile(me.getFilePath())
}

func (me *database) openFile(filePath string) *sql.DB {
	const journalMode = "_journal_mode=WAL"
	var busyTimeout = "_busy_timeout=" + strconv.Itoa(int(dbTimeout))
	var url = "file:" + filePath + "?" + journalMode + "&" + busyTimeout
	return assertResultError(sql.Open("sqlite3", url))
}

func (me *database) close(db *sql.DB) *sql.DB {
	if db != nil {
		assertError(db.Close())
	}
	return nil
}

func (me *database) getFilePath() string {
	return me.dataDirectory + "/hinst-website.db"
}

func (me *database) collectGarbage() {
	var db = me.open()
	defer me.close(db)

	assertResultError(db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"))
	assertResultError(db.Exec("VACUUM;"))
	assertResultError(db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"))
}
