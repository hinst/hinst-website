package main

import (
	"database/sql"
	"fmt"
	"os"
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

func (me *database) getGoalIds() (goalIds []int) {
	var files = assertResultError(os.ReadDir(me.dataDirectory))
	for _, file := range files {
		if file.IsDir() && goalIdStringMatcher.MatchString(file.Name()) {
			var goalId = assertResultError(strconv.Atoi(file.Name()))
			goalIds = append(goalIds, goalId)
		}
	}
	return
}

func (me *database) migrate() {
	// merge old and new database formats
	var newDb = me.openFile("C:\\Dev\\SmartProgress-or\\downloader\\data\\hinst-website.db")
	defer me.close(newDb)
	var oldDb = me.open()
	defer me.close(oldDb)
	me.forEachGoalPost(func(row *goalPostRow) {
		var execResult = assertResultError(
			newDb.Exec("UPDATE goalPosts SET isPublic = ? WHERE goalId = ? AND dateTime = ?",
				row.isPublic, row.goalId, row.dateTime.UTC().Unix()))
		fmt.Println(row.String())
		fmt.Printf("%v\n", assertResultError(execResult.RowsAffected()))
	})
}

func (me *database) getGoalPost(goalId int, dateTime time.Time) (result *goalPostRow) {
	var db = me.open()
	defer me.close(db)
	var rows = assertResultError(
		db.Query("SELECT * FROM goalPosts WHERE goalId = ? AND dateTime = ?",
			goalId, dateTime.UTC().Unix()),
	)
	if rows.Next() {
		assertError(rows.Err())
		result = new(goalPostRow)
		result.scan(rows)
	}
	return
}

func (me *database) getPostsByDates(ids []int64) (results map[time.Time]goalPostRow) {
	var db = me.open()
	defer me.close(db)
	var queryText = "SELECT * FROM goalPosts WHERE dateTime IN (" + convertInt64ArrayToSequelString(ids) + ")"
	var rows = assertResultError(db.Query(queryText))
	results = make(map[time.Time]goalPostRow)
	for rows.Next() {
		var row goalPostRow
		row.scan(rows)
		results[row.dateTime.UTC()] = row
	}
	return
}

func (me *database) getGoalPostVisibilities(goalId int) (result map[time.Time]bool) {
	var db = me.open()
	defer me.close(db)
	var rows = assertResultError(
		db.Query("SELECT dateTime, isPublic FROM goalPosts WHERE goalId = ?", goalId),
	)
	result = make(map[time.Time]bool)
	for rows.Next() {
		var dateTimeMilliseconds int64
		var isPublic bool
		assertError(rows.Scan(&dateTimeMilliseconds, &isPublic))
		result[time.Unix(dateTimeMilliseconds, 0).UTC()] = isPublic
	}
	return
}

func (me *database) setGoalPostPublic(row *goalPostRow) {
	var db = me.open()
	defer me.close(db)
	assertResultError(
		db.Exec("INSERT INTO goalPosts (goalId, dateTime, isPublic) VALUES (?, ?, ?) "+
			"ON CONFLICT(goalId, dateTime) DO UPDATE SET isPublic = ?",
			row.goalId, row.dateTime.UTC().Unix(), row.isPublic, row.isPublic),
	)
}

func (me *database) forEachGoalPost(callback func(row *goalPostRow)) {
	var db = me.open()
	defer me.close(db)
	var rows = assertResultError(db.Query("SELECT * FROM goalPosts"))
	for rows.Next() {
		var row goalPostRow
		row.scan(rows)
		callback(&row)
	}
	assertError(rows.Err())
}

func (me *database) collectGarbage() {
	var db = me.open()
	defer me.close(db)

	assertResultError(db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"))
	assertResultError(db.Exec("VACUUM;"))
	assertResultError(db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"))
}
