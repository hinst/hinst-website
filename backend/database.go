package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
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
	const journalMode = "_journal_mode=WAL"
	var busyTimeout = "_busy_timeout=" + strconv.Itoa(int(dbTimeout))
	var url = "file:" + me.getFilePath() + "?" + journalMode + "&" + busyTimeout
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
	var goalIds = me.getGoalIds()
	var db = me.open()
	defer me.close(db)
	for _, goalId := range goalIds {
		var goalDirectory = filepath.Join(me.dataDirectory, strconv.Itoa(goalId))
		var publicPostsPath = goalDirectory + "/" + publicPostsFileName
		if checkFileExists(publicPostsPath) {
			var publicPostsText = readTextFile(publicPostsPath)
			var publicPosts = strings.Split(publicPostsText, "\n")
			for i := range publicPosts {
				publicPosts[i] = strings.TrimSpace(publicPosts[i])
			}
			var fileNames = getGoalFiles(goalDirectory)
			for _, fileName := range fileNames {
				var dateTime = assertResultError(parseStoredGoalFileDate(getFileNameWithoutExtension(fileName)))
				var dateTimeText = dateTime.Format(smartProgressTimeFormat)
				var isPublic = slices.Contains(publicPosts, dateTimeText)
				assertResultError(
					db.Exec("INSERT INTO goalPosts (goalId, dateTime, isPublic) VALUES (?, ?, ?)",
						goalId, dateTime.UTC().Unix(), isPublic))
			}
		}
	}
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

func (me *database) getGoalPostRows(goalId int) (goalPostRows []goalPostRow) {
	var db = me.open()
	defer me.close(db)
	var rows = assertResultError(
		db.Query("SELECT * FROM goalPosts WHERE isPublic = 1 AND goalId = ?", goalId),
	)
	for rows.Next() {
		var row goalPostRow
		row.scan(rows)
		goalPostRows = append(goalPostRows, row)
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

func (me *database) collectGarbage() {
	var db = me.open()
	defer me.close(db)

	assertResultError(db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"))
	assertResultError(db.Exec("VACUUM;"))
	assertResultError(db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"))
}
