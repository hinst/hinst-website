package main

import (
	"database/sql"
	"strconv"
	"time"
)

// Milliseconds
var DB_TIMEOUT = time.Hour / time.Millisecond

type Database struct {
	filePath string
}

func (me *Database) init(filePath string) {
	me.filePath = filePath
}

func (me *Database) open() *sql.DB {
	return assertResultError(sql.Open("sqlite3", "file:"+me.filePath+
		"?_journal_mode=WAL&_busy_timeout="+strconv.Itoa(int(DB_TIMEOUT))))
}

func (me *Database) close(db *sql.DB) *sql.DB {
	if db != nil {
		assertError(db.Close())
	}
	return nil
}
