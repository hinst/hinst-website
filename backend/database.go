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
var DB_TIMEOUT = time.Hour / time.Millisecond

//go:embed schema.sql
var DB_SCHEMA string

type Database struct {
	dataDirectory string
}

func (me *Database) init(dataDirectory string) {
	me.dataDirectory = dataDirectory
	var db = me.open()
	defer me.close(db)
	assertError(db.Ping())
	assertResultError(db.Exec(DB_SCHEMA))
	me.collectGarbage()
}

func (me *Database) open() *sql.DB {
	const journalMode = "_journal_mode=WAL"
	var busyTimeout = "_busy_timeout=" + strconv.Itoa(int(DB_TIMEOUT))
	return assertResultError(sql.Open(
		"sqlite3",
		"file:"+me.getFilePath()+
			"?"+journalMode+
			"&"+busyTimeout,
	))
}

func (me *Database) close(db *sql.DB) *sql.DB {
	if db != nil {
		assertError(db.Close())
	}
	return nil
}

func (me *Database) getFilePath() string {
	return me.dataDirectory + "/hinst-website.db"
}

func (me *Database) getGoalIds() (goalIds []int) {
	var files = assertResultError(os.ReadDir(me.dataDirectory))
	for _, file := range files {
		if file.IsDir() && goalIdStringMatcher.MatchString(file.Name()) {
			var goalId = assertResultError(strconv.Atoi(file.Name()))
			goalIds = append(goalIds, goalId)
		}
	}
	return
}

func (me *Database) migrate() {
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
					db.Exec("INSERT INTO goalPosts (goalId, dateTime, isPublic) VALUES (?, ?, ?) RETURNING rowid",
						goalId, dateTime.UTC().Unix(), isPublic))
			}
		}
	}
}

func (me *Database) collectGarbage() {
	var db = me.open()
	defer me.close(db)

	assertResultError(db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"))
	assertResultError(db.Exec("VACUUM;"))
	assertResultError(db.Exec("PRAGMA wal_checkpoint(TRUNCATE);"))
}
