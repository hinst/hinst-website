package main

import (
	"slices"
	"time"

	"golang.org/x/text/language"
)

func (me *database) setGoalPostPublic(row goalPostRow) {
	var db = me.open()
	defer me.close(db)
	assertResultError(
		db.Exec("UPDATE goalPosts SET isPublic = ? WHERE goalId = ? AND dateTime = ?",
			row.isPublic, row.goalId, row.dateTime.UTC().Unix()),
	)
}

func (me *database) setGoalPostText(goalId int64, dateTime time.Time, supportedLanguage language.Tag, text string) {
	var db = me.open()
	defer me.close(db)
	assertCondition(slices.Contains(supportedLanguages, supportedLanguage),
		func() string { return "Unsupported language: " + supportedLanguage.String() })
	var languageName = ""
	if supportedLanguage != supportedLanguages[0] {
		languageName = getLanguageName(supportedLanguage)
	}
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

func (me *database) getGoalPost(goalId int64, dateTime time.Time) (result *goalPostRow) {
	var db = me.open()
	defer me.close(db)
	var queryText = "SELECT * FROM goalPosts WHERE goalId = ? AND dateTime = ?"
	var rows = assertResultError(db.Query(queryText, goalId, dateTime.UTC().Unix()))
	if rows.Next() {
		result = new(goalPostRow)
		result.scan(rows)
	}
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
