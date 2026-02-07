package server

import (
	"context"
	"log"
	"slices"
	"strings"
	"time"

	"golang.org/x/text/language"
)

func (me *database) setGoalPostPublic(row goalPostRow) int64 {
	var query = "UPDATE goalPosts SET isPublic = $1 WHERE goalId = $2 AND dateTime = $3"
	var result = assertResultError(me.pool.Exec(context.Background(), query,
		row.isPublic, row.goalId, row.getDateTime().UTC().Unix()))
	return result.RowsAffected()
}

func (me *database) setGoalPostText(goalId int64, dateTime time.Time, supportedLanguage language.Tag, text string) int64 {
	var textField = "text" + me.getLanguagePostfix(supportedLanguage)
	var queryText = "UPDATE goalPosts SET " + textField + " = $1 WHERE goalId = $2 AND dateTime = $3"
	var dateTimeEpoch = dateTime.UTC().Unix()
	var result = assertResultError(me.pool.Exec(context.Background(), queryText, text, goalId, dateTimeEpoch))
	return result.RowsAffected()
}

func (me *database) setGoalPostTitle(goalId int64, dateTime time.Time, supportedLanguage language.Tag, text string) int64 {
	var titleField = "title" + me.getLanguagePostfix(supportedLanguage)
	var queryText = "UPDATE goalPosts SET " + titleField + " = $1 WHERE goalId = $2 AND dateTime = $3"
	var dateTimeEpoch = dateTime.UTC().Unix()
	var result = assertResultError(me.pool.Exec(context.Background(), queryText, text, goalId, dateTimeEpoch))
	return result.RowsAffected()
}

func (me *database) forEachGoalPost(callback func(row *goalPostRow) bool, selector string, sortByDate int) {
	var querySql = "SELECT " + selector + " FROM goalPosts"
	if sortByDate != 0 {
		querySql += " ORDER BY dateTime " + ifElse(sortByDate > 0, "ASC", "DESC")
	}
	var rows = assertResultError(me.pool.Query(context.Background(), querySql))
	defer rows.Close()
	for rows.Next() {
		var row goalPostRow
		row.scan(rows)
		assertError(rows.Err())
		if !callback(&row) {
			break
		}
	}
}

func (me *database) getGoals() (results []goalRecord) {
	var rows = assertResultError(me.pool.Query(context.Background(), "SELECT id, title FROM goals"))
	defer rows.Close()
	for rows.Next() {
		var record goalRecord
		assertError(rows.Scan(&record.Id, &record.Title))
		results = append(results, record)
	}
	return
}

func (me *database) getGoal(goalId int64) (result *goalRecord) {
	var queryText = "SELECT id, title FROM goals WHERE id = $1"
	var rows = assertResultError(me.pool.Query(context.Background(), queryText, goalId))
	defer rows.Close()
	if rows.Next() {
		result = new(goalRecord)
		assertError(rows.Scan(&result.Id, &result.Title))
	}
	return
}

func (me *database) getGoalPost(goalId int64, dateTime time.Time) (result *goalPostRow) {
	var queryText = "SELECT * FROM goalPosts WHERE goalId = $1 AND dateTime = $2"
	var rows = assertResultError(me.pool.Query(context.Background(), queryText, goalId, dateTime.UTC().Unix()))
	defer rows.Close()
	if rows.Next() {
		result = new(goalPostRow)
		result.scan(rows)
	}
	return
}

func (me *database) getGoalPostImage(goalId int64, dateTime time.Time, index int) (result *goalPostImageRow) {
	var queryText = "SELECT contentType, file FROM goalPostImages" +
		" WHERE goalId = $1 AND parentDateTime = $2 AND sequenceIndex = $3"
	var rows = assertResultError(me.pool.Query(context.Background(), queryText, goalId, dateTime.UTC().Unix(), index))
	defer rows.Close()
	if rows.Next() {
		result = new(goalPostImageRow)
		assertError(rows.Scan(&result.contentType, &result.file))
	}
	return
}

func (me *database) getGoalPostImageCount(goalId int64, dateTime time.Time) (count int) {
	var queryText = "SELECT COUNT(*) FROM goalPostImages WHERE goalId = $1 AND parentDateTime = $2"
	var row = me.pool.QueryRow(context.Background(), queryText, goalId, dateTime.UTC().Unix())
	assertError(row.Scan(&count))
	return
}

func (me *database) getGoalPosts(goalId int64, includePrivate bool, language language.Tag) (results []goalPostRecord) {
	var titleField = "title" + me.getLanguagePostfix(language)
	var queryText = "SELECT goalId, dateTime, isPublic, type, " + titleField + " FROM goalPosts WHERE goalId = $1"
	if !includePrivate {
		queryText += " AND isPublic = 1"
	}
	var rows = assertResultError(me.pool.Query(context.Background(), queryText, goalId))
	defer rows.Close()
	for rows.Next() {
		var record goalPostRecord
		assertError(rows.Scan(&record.GoalId, &record.DateTime, &record.IsPublic, &record.Type, &record.Title))
		results = append(results, record)
	}
	return
}

func (database) getLanguagePostfix(supportedLanguage language.Tag) string {
	assertCondition(slices.Contains(supportedLanguages, supportedLanguage),
		func() string { return "Unsupported language: " + supportedLanguage.String() })
	var languageName = ""
	if supportedLanguage != supportedLanguages[0] {
		languageName = getLanguageName(supportedLanguage)
	}
	return languageName
}

func (me *database) searchGoalPosts(
	queryText string, supportedLanguage language.Tag, includePrivate bool, limit int,
) (results []*goalPostRow) {
	queryText = strings.ToUpper(queryText)
	queryText = normalizeString(queryText)
	me.forEachGoalPost(func(row *goalPostRow) bool {
		if limit <= len(results) {
			return false
		}
		var isVisible = includePrivate || row.isPublic
		if !isVisible {
			return true
		}
		var title = strings.ToUpper(row.getTranslatedTitle(supportedLanguage))
		var text = stripHtml(strings.ToUpper(row.getTranslatedText(supportedLanguage)))
		if strings.Contains(title, queryText) || strings.Contains(text, queryText) {
			results = append(results, row)
		}
		return true
	}, (goalPostRow{}).getSelectorForLanguage(supportedLanguage), -1)
	return
}

func (me *database) migrate() {
	log.Println("Migrating table urlPings")
	assertResultError(me.pool.Exec(context.Background(), "ALTER TABLE urlPings ADD COLUMN googlePingedManuallyAt BIGINT"))
}
