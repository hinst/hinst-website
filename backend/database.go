package main

import (
	"database/sql"
	"os"
	"path/filepath"
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
}

func (me *Database) open() *sql.DB {
	return assertResultError(sql.Open("sqlite3", "file:"+me.getFilePath()+
		"?_journal_mode=WAL&_busy_timeout="+strconv.Itoa(int(DB_TIMEOUT))))
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
			for _, publicPostDateTimeText := range publicPosts {
				publicPostDateTimeText = strings.TrimSpace(publicPostDateTimeText)
				var publicPostDateTime = assertResultError(parseSmartProgressDate(publicPostDateTimeText))
				publicPostDateTime = publicPostDateTime
			}
		}
	}
}
