package server

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/base"
	"github.com/hinst/hinst-website/server/db_objects"
	"github.com/hinst/hinst-website/server/rest_objects"
	"golang.org/x/text/language"
)

func (me *database) setGoalPostPublic(row db_objects.GoalPostRow) int64 {
	var query = "UPDATE " + row.GetTableName() + " SET isPublic = $1 WHERE goalId = $2 AND dateTime = $3"
	var result = gophers.AssertResultError(me.pool.Exec(context.Background(), query,
		row.IsPublic, row.GoalId, row.GetDateTime().UTC().Unix()))
	return result.RowsAffected()
}

func (me *database) setGoalPostSearchIndexingEnabled(row db_objects.GoalPostRow) int64 {
	var query = "UPDATE " + row.GetTableName() + " SET searchIndexingEnabled = $1 WHERE goalId = $2 AND dateTime = $3"
	var result = gophers.AssertResultError(me.pool.Exec(context.Background(), query,
		row.SearchIndexingEnabled, row.GoalId, row.GetDateTime().UTC().Unix()))
	return result.RowsAffected()
}

func (me *database) setGoalPostSearchIndexingStatus(
	goalId int64, dateTime time.Time, status string, checkedAt time.Time,
) int64 {
	var queryText = "UPDATE " + (db_objects.GoalPostRow{}).GetTableName() +
		" SET googleSearchIndexingStatus = $1, googleSearchIndexingStatusCheckedAt = $2" +
		" WHERE goalId = $3 AND dateTime = $4"
	var result = gophers.AssertResultError(me.pool.Exec(context.Background(), queryText,
		status, checkedAt.UTC().Unix(), goalId, dateTime.UTC().Unix()))
	return result.RowsAffected()
}

func (me *database) setGoalPostGooglePingedAt(goalId int64, dateTime time.Time, googlePingedAt time.Time) int64 {
	var queryText = "UPDATE " + (db_objects.GoalPostRow{}).GetTableName() +
		" SET googlePingedAt = $1 WHERE goalId = $2 AND dateTime = $3"
	var result = gophers.AssertResultError(me.pool.Exec(context.Background(), queryText,
		googlePingedAt.UTC().Unix(), goalId, dateTime.UTC().Unix()))
	return result.RowsAffected()
}

// Sets the element of the text array that corresponds to the given language
func (me *database) setGoalPostText(goalId int64, dateTime time.Time, supportedLanguage language.Tag, text string) int64 {
	return me.setGoalPostLanguageArrayElement("text", goalId, dateTime, supportedLanguage, text)
}

func (me *database) setGoalPostTitle(goalId int64, dateTime time.Time, supportedLanguage language.Tag, text string) int64 {
	return me.setGoalPostLanguageArrayElement("title", goalId, dateTime, supportedLanguage, text)
}

// columnName is either "text" or "title"; the array elements are aligned with base.SupportedLanguages
func (me *database) setGoalPostLanguageArrayElement(
	columnName string, goalId int64, dateTime time.Time, supportedLanguage language.Tag, text string,
) int64 {
	var languageIndex = base.GetLanguageIndex(supportedLanguage)
	var ctx = context.Background()
	var tx = gophers.AssertResultError(me.pool.Begin(ctx))
	defer tx.Rollback(ctx) // Safe to ignore the result; Rollback after Commit is a no-op
	var tableName = (db_objects.GoalPostRow{}).GetTableName()
	var dateTimeEpoch = dateTime.UTC().Unix()
	var current []string
	gophers.AssertError(tx.QueryRow(ctx,
		"SELECT "+columnName+" FROM "+tableName+" WHERE goalId = $1 AND dateTime = $2",
		goalId, dateTimeEpoch).Scan(&current))
	for len(current) <= languageIndex {
		current = append(current, "")
	}
	current[languageIndex] = text
	var result = gophers.AssertResultError(tx.Exec(ctx,
		"UPDATE "+tableName+" SET "+columnName+" = $1 WHERE goalId = $2 AND dateTime = $3",
		current, goalId, dateTimeEpoch))
	gophers.AssertError(tx.Commit(ctx))
	return result.RowsAffected()
}

// Callback should return true to continue the loop, return false to break the loop early.
func (me *database) forEachGoalPost(callback func(row *db_objects.GoalPostRow) bool, selector string, sortByDate int) {
	var tableName = (db_objects.GoalPostRow{}).GetTableName()
	var querySql = "SELECT " + selector + " FROM " + tableName
	if sortByDate != 0 {
		querySql += " ORDER BY dateTime " + gophers.IfElse(sortByDate > 0, "ASC", "DESC")
	}
	var rows = gophers.AssertResultError(me.pool.Query(context.Background(), querySql))
	defer rows.Close()
	for rows.Next() {
		var row = new(db_objects.GoalPostRow)
		row.Scan(rows)
		gophers.AssertError(rows.Err())
		if !callback(row) {
			break
		}
	}
}

func (me *database) getGoals() (results []db_objects.GoalRow) {
	var tableName = (db_objects.GoalRow{}).GetTableName()
	var fields = strings.Join(gophers.GetFieldNames[db_objects.GoalRow](), ",")
	var rows = gophers.AssertResultError(me.pool.Query(context.Background(), "SELECT "+fields+" FROM "+tableName+" ORDER BY id"))
	defer rows.Close()
	for rows.Next() {
		var record db_objects.GoalRow
		record.Scan(rows)
		results = append(results, record)
	}
	return
}

func (me *database) getGoal(goalId int64) (result *db_objects.GoalRow) {
	var tableName = (db_objects.GoalRow{}).GetTableName()
	var fields = strings.Join(gophers.GetFieldNames[db_objects.GoalRow](), ",")
	var queryText = "SELECT " + fields + " FROM " + tableName + " WHERE id = $1"
	var rows = gophers.AssertResultError(me.pool.Query(context.Background(), queryText, goalId))
	defer rows.Close()
	if rows.Next() {
		result = new(db_objects.GoalRow)
		result.Scan(rows)
	}
	return
}

func (me *database) getGoalImage(goalId int64) (imageData []byte, imageContentType string) {
	var tableName = (db_objects.GoalRow{}).GetTableName()
	var queryText = "SELECT imageData, imageContentType FROM " + tableName + " WHERE id = $1"
	var rows = gophers.AssertResultError(me.pool.Query(context.Background(), queryText, goalId))
	defer rows.Close()
	if rows.Next() {
		gophers.AssertError(rows.Scan(&imageData, &imageContentType))
	}
	return
}

func (me *database) insertGoal(row *db_objects.GoalRow) (isInserted bool) {
	var columnNames = row.GetAllColumns()
	var query = "INSERT INTO " + row.GetTableName() +
		" (" + strings.Join(columnNames, ",") + ")" +
		" VALUES (" + me.buildPlaceholders(len(columnNames)) + ")" +
		" ON CONFLICT DO NOTHING"
	var result = gophers.AssertResultError(me.pool.Exec(context.Background(), query,
		db_objects.GetAllColumnValues(row)...))
	return result.RowsAffected() == 1
}

func (me *database) updateGoalSmart(row *db_objects.GoalRow) (isUpdated bool) {
	var columnNames = row.GetSmartColumns()
	var values = gophers.AssertResultError(gophers.GetFieldValuesByNames(row, columnNames))
	values = append(values, row.Id)
	var assignments = make([]string, len(columnNames))
	for i, columnName := range columnNames {
		assignments[i] = columnName + " = $" + strconv.Itoa(i+1)
	}
	var idPlaceholder = "$" + strconv.Itoa(len(columnNames)+1)
	var query = "UPDATE " + row.GetTableName() +
		" SET " + strings.Join(assignments, ",") +
		" WHERE id = " + idPlaceholder
	var result = gophers.AssertResultError(me.pool.Exec(context.Background(), query, values...))
	return result.RowsAffected() == 1
}

// Returns true if the goal post was inserted;
// Returns false if goal post with this key already exists;
func (me *database) insertGoalPost(row *db_objects.GoalPostRow) (isInserted bool) {
	var columnNames = row.GetAllColumns()
	var query = "INSERT INTO " + row.GetTableName() +
		" (" + strings.Join(columnNames, ",") + ")" +
		" VALUES (" + me.buildPlaceholders(len(columnNames)) + ")" +
		" ON CONFLICT DO NOTHING"
	var result = gophers.AssertResultError(me.pool.Exec(context.Background(), query,
		db_objects.GetAllColumnValues(row)...))
	return result.RowsAffected() == 1
}

func (me *database) saveGoalPostComment(row db_objects.GoalPostCommentRow) {
	var columnNames = row.GetAllColumns()
	var query = "INSERT INTO " + row.GetTableName() +
		" (" + strings.Join(columnNames, ",") + ")" +
		" VALUES (" + me.buildPlaceholders(len(columnNames)) + ")" +
		" ON CONFLICT (goalId, parentDateTime, dateTime, smartProgressUserId)" +
		" DO UPDATE SET username = excluded.username, text = excluded.text"
	gophers.AssertResultError(me.pool.Exec(context.Background(), query, db_objects.GetAllColumnValues(&row)...))
}

func (me *database) saveGoalPostImage(row db_objects.GoalPostImageRow) {
	var columnNames = row.GetAllColumns()
	var query = "INSERT INTO " + row.GetTableName() +
		" (" + strings.Join(columnNames, ",") + ")" +
		" VALUES (" + me.buildPlaceholders(len(columnNames)) + ")" +
		" ON CONFLICT (goalId, parentDateTime, sequenceIndex)" +
		" DO UPDATE SET contentType = excluded.contentType, file = excluded.file"
	gophers.AssertResultError(me.pool.Exec(context.Background(), query, db_objects.GetAllColumnValues(&row)...))
}

func (me *database) getGoalPost(goalId int64, dateTime time.Time) (result *db_objects.GoalPostRow) {
	var tableName = (db_objects.GoalPostRow{}).GetTableName()
	var fields = db_objects.GoalPostRow{}.GetAllFieldSelector()
	var queryText = "SELECT " + fields + " FROM " + tableName + " WHERE goalId = $1 AND dateTime = $2"
	var rows = gophers.AssertResultError(me.pool.Query(context.Background(), queryText, goalId, dateTime.UTC().Unix()))
	defer rows.Close()
	if rows.Next() {
		result = new(db_objects.GoalPostRow)
		result.Scan(rows)
	}
	return
}

func (me *database) getGoalPostImage(goalId int64, dateTime time.Time, index int) (result *db_objects.GoalPostImageRow) {
	var tableName = (db_objects.GoalPostImageRow{}).GetTableName()
	var queryText = "SELECT contentType, file FROM " + tableName +
		" WHERE goalId = $1 AND parentDateTime = $2 AND sequenceIndex = $3"
	var rows = gophers.AssertResultError(me.pool.Query(context.Background(), queryText, goalId, dateTime.UTC().Unix(), index))
	defer rows.Close()
	if rows.Next() {
		result = new(db_objects.GoalPostImageRow)
		gophers.AssertError(rows.Scan(&result.ContentType, &result.File))
	}
	return
}

func (me *database) getGoalPostImageCount(goalId int64, dateTime time.Time) (count int) {
	var tableName = (db_objects.GoalPostImageRow{}).GetTableName()
	var queryText = "SELECT COUNT(*) FROM " + tableName + " WHERE goalId = $1 AND parentDateTime = $2"
	var row = me.pool.QueryRow(context.Background(), queryText, goalId, dateTime.UTC().Unix())
	gophers.AssertError(row.Scan(&count))
	return
}

func (me *database) getGoalPosts(goalId int64, includePrivate bool, language language.Tag) (results []rest_objects.GoalPostHeader) {
	var tableName = (db_objects.GoalPostRow{}).GetTableName()
	var queryText = "SELECT goalId, dateTime, isPublic, type, title FROM " + tableName + " WHERE goalId = $1"
	if !includePrivate {
		queryText += " AND isPublic = TRUE"
	}
	queryText += " ORDER BY dateTime DESC"
	var languageIndex = base.GetLanguageIndex(language)
	var rows = gophers.AssertResultError(me.pool.Query(context.Background(), queryText, goalId))
	defer rows.Close()
	for rows.Next() {
		var record rest_objects.GoalPostHeader
		var title []string
		gophers.AssertError(rows.Scan(&record.GoalId, &record.DateTime, &record.IsPublic, &record.Type, &title))
		if languageIndex < len(title) {
			record.Title = title[languageIndex]
		}
		results = append(results, record)
	}
	return
}

func (me *database) searchGoalPosts(
	queryText string, supportedLanguage language.Tag, includePrivate bool, limit int,
) (results []*db_objects.GoalPostRow) {
	queryText = strings.ToUpper(queryText)
	queryText = gophers.NormalizeString(queryText)
	me.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		if limit <= len(results) {
			return false
		}
		var isVisible = includePrivate || row.IsPublic
		if !isVisible {
			return true
		}
		var title = strings.ToUpper(row.GetTranslatedTitle(supportedLanguage))
		var text = strings.ToUpper(row.GetTranslatedText(supportedLanguage))
		if strings.Contains(title, queryText) || strings.Contains(text, queryText) {
			results = append(results, row)
		}
		return true
	}, (db_objects.GoalPostRow{}).GetAllFieldSelector(), -1)
	return
}
