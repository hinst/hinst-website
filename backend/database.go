package main

import (
	"database/sql"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/language"
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
	var queryText = "SELECT isPublic, text, textEnglish, textGerman, type FROM goalPosts WHERE goalId = ? AND dateTime = ?"
	var row = db.QueryRow(queryText, goalId, dateTime.UTC().Unix())
	assertError(row.Err())
	result.goalId = goalId
	result.dateTime = dateTime
	assertError(row.Scan(&result.isPublic, &result.text, &result.textEnglish, &result.textGerman, &result.typeString))
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

func (me *database) getGoalPosts(goalId int64, includePrivate bool) (results []goalPostRecord) {
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
	var oldDb = me.open()
	defer me.close(oldDb)
	me.forEachGoalPost(func(row *goalPostRow) {
		var translatedFilesFolder = filepath.Join(me.dataDirectory, getStringFromInt64(row.goalId), "translated")
		var translatedFileNames = assertResultError(os.ReadDir(translatedFilesFolder))
		var minDiff int64 = math.MaxInt64
		var matchedDateTimeText = ""
		for _, translatedFileName := range translatedFileNames {
			var dateTimeText = translatedFileName.Name()[0:len("2025-01-01_00-00-00")]
			var dateTimeForFile = assertResultError(time.Parse("2006-01-02_15-04-05", dateTimeText))
			var diff = absInt64(row.dateTime.Unix() - dateTimeForFile.Unix())
			if diff < minDiff && diff < 60*60*12 {
				minDiff = diff
				matchedDateTimeText = dateTimeText
			}
		}
		if matchedDateTimeText == "" {
			log.Printf("⚠️ Cannot find translated file for row goalId=%v dateTime=%v\n", row.goalId, row.dateTime)
			return
		}
		for _, supportedLanguage := range supportedLanguages[1:] {
			var filePath = filepath.Join(translatedFilesFolder, matchedDateTimeText+"."+supportedLanguage.String()+".html")
			if !checkFileExists(filePath) {
				log.Printf("⚠️ Cannot find translated file for row goalId=%v dateTime=%v language=%v\n",
					row.goalId, row.dateTime, supportedLanguage.String())
				continue
			}
			var fileText = readTextFile(filePath)
			me.setTranslatedText(row.goalId, row.dateTime, supportedLanguage, fileText)
		}
	})
}

func (me *database) setGoalPostPublic(row goalPostRow) {
	var db = me.open()
	defer me.close(db)
	assertResultError(
		db.Exec("INSERT INTO goalPosts (goalId, dateTime, isPublic) VALUES (?, ?, ?) "+
			"ON CONFLICT(goalId, dateTime) DO UPDATE SET isPublic = ?",
			row.goalId, row.dateTime.UTC().Unix(), row.isPublic, row.isPublic),
	)
}

func (me *database) setTranslatedText(goalId int64, dateTime time.Time, supportedLanguage language.Tag, text string) {
	var db = me.open()
	defer me.close(db)
	assertCondition(slices.Contains(supportedLanguages[1:], supportedLanguage),
		func() string { return "Unsupported language: " + supportedLanguage.String() })
	var languageName = getLanguageName(supportedLanguage)
	var queryText = "UPDATE goalPosts SET text" + languageName + " = ? WHERE goalId = ? AND dateTime = ?"
	var dateTimeEpoch = dateTime.UTC().Unix()
	var result = assertResultError(db.Exec(queryText, text, goalId, dateTimeEpoch))
	var changedRowCount = assertResultError(result.RowsAffected())
	assertCondition(changedRowCount > 0,
		func() string {
			return "Cannot update translated text for goalId=" +
				getStringFromInt64(goalId) + " dateTime=" + getStringFromInt64(dateTimeEpoch)
		})
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
