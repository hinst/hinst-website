package main

import (
	"database/sql"
	"fmt"
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

func (me *database) getGoals() (results []goalRecord) {
	var db = me.open()
	defer me.close(db)
	var rows = assertResultError(db.Query("SELECT id, title FROM goals"))
	for rows.Next() {
		var record goalRecord
		assertError(rows.Scan(&record.Id, &record.Title))
		results = append(results, record)
	}
	return
}

func (me *database) getGoal(goalId int64) (result goalRecord) {
	var db = me.open()
	defer me.close(db)
	var row = db.QueryRow("SELECT id, title FROM goals WHERE id = ?", goalId)
	assertError(row.Err())
	assertError(row.Scan(&result.Id, &result.Title))
	return
}

func (me *database) getGoalPost(goalId int64, dateTime time.Time) (result goalPostRow) {
	var db = me.open()
	defer me.close(db)
	var row = db.QueryRow("SELECT isPublic, text FROM goalPosts WHERE goalId = ? AND dateTime = ?",
		goalId, dateTime.UTC().Unix())
	assertError(row.Err())
	result.GoalId = goalId
	result.DateTime = dateTime
	assertError(row.Scan(&result.IsPublic, &result.Text))
	return
}

func (me *database) getGoalPostImages(goalId int64, dateTime time.Time) (results []goalPostImageRow) {
	var db = me.open()
	defer me.close(db)
	var queryText = "SELECT contentType, file FROM goalPostImages" +
		" WHERE goalId = ? AND parentDateTime = ?" +
		" ORDER BY sequenceIndex"
	var rows = assertResultError(db.Query(queryText, goalId, dateTime.UTC().Unix()))
	for rows.Next() {
		var record goalPostImageRow
		assertError(rows.Scan(&record.contentType, &record.file))
		results = append(results, record)
	}
	return
}

func (me *database) getGoalPosts(goalId int, includePrivate bool) (results []goalPostRecord) {
	var db = me.open()
	defer me.close(db)
	var queryText = "SELECT goalId, dateTime, isPublic, type FROM goalPosts WHERE goalId = ?"
	if !includePrivate {
		queryText += " AND isPublic = 1"
	}
	var rows = assertResultError(db.Query(queryText, goalId))
	for rows.Next() {
		var record goalPostRecord
		assertError(rows.Scan(&record.GoalId, &record.DateTime, &record.IsPublic, &record.Type))
		results = append(results, record)
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
		var matchedRow = newDb.QueryRow("SELECT dateTime FROM goalPosts WHERE goalId = ? ORDER BY abs(dateTime - ?) LIMIT 1",
			row.GoalId, row.DateTime.UTC().Unix())
		assertError(matchedRow.Err())
		var matchedDateTime int64
		assertError(matchedRow.Scan(&matchedDateTime))
		assertCondition(matchedDateTime != 0, func() string { return "Cannot find matching date time" })
		var difference time.Duration = time.Duration(absInt64(row.DateTime.UTC().Unix()-matchedDateTime)) * time.Second
		fmt.Printf("Difference: %v\n", difference)
		var execResult = assertResultError(
			newDb.Exec("UPDATE goalPosts SET isPublic = ? WHERE goalId = ? AND dateTime = ?",
				row.IsPublic, row.GoalId, matchedDateTime))
		fmt.Println(row.String())
		fmt.Printf("%v\n", assertResultError(execResult.RowsAffected()))
	})
}

func (me *database) setGoalPostPublic(row *goalPostRow) {
	var db = me.open()
	defer me.close(db)
	assertResultError(
		db.Exec("INSERT INTO goalPosts (goalId, dateTime, isPublic) VALUES (?, ?, ?) "+
			"ON CONFLICT(goalId, dateTime) DO UPDATE SET isPublic = ?",
			row.GoalId, row.DateTime.UTC().Unix(), row.IsPublic, row.IsPublic),
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
