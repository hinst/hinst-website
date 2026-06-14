package server

import (
	"context"
	"strings"
	"time"

	"github.com/hinst/go-common"
	"github.com/hinst/hinst-website/server/db_objects"
	"github.com/hinst/hinst-website/server/rest_objects"
	"golang.org/x/text/language"
)

func (me *database) setGoalPostPublic(row db_objects.GoalPostRow) int64 {
	var query = "UPDATE goalPosts SET isPublic = $1 WHERE goalId = $2 AND dateTime = $3"
	var result = common.AssertResultError(me.pool.Exec(context.Background(), query,
		row.IsPublic, row.GoalId, row.GetDateTime().UTC().Unix()))
	return result.RowsAffected()
}

func (me *database) setGoalPostText(goalId int64, dateTime time.Time, supportedLanguage language.Tag, text *string) int64 {
	var textField = "text" + db_objects.GetLanguagePostfix(supportedLanguage)
	var queryText = "UPDATE goalPosts SET " + textField + " = $1 WHERE goalId = $2 AND dateTime = $3"
	var dateTimeEpoch = dateTime.UTC().Unix()
	var result = common.AssertResultError(me.pool.Exec(context.Background(), queryText, text, goalId, dateTimeEpoch))
	return result.RowsAffected()
}

func (me *database) setGoalPostTitle(goalId int64, dateTime time.Time, supportedLanguage language.Tag, text string) int64 {
	var titleField = "title" + db_objects.GetLanguagePostfix(supportedLanguage)
	var queryText = "UPDATE goalPosts SET " + titleField + " = $1 WHERE goalId = $2 AND dateTime = $3"
	var dateTimeEpoch = dateTime.UTC().Unix()
	var result = common.AssertResultError(me.pool.Exec(context.Background(), queryText, text, goalId, dateTimeEpoch))
	return result.RowsAffected()
}

func (me *database) forEachGoalPost(callback func(row *db_objects.GoalPostRow) bool, selector string, sortByDate int) {
	var querySql = "SELECT " + selector + " FROM goalPosts"
	if sortByDate != 0 {
		querySql += " ORDER BY dateTime " + common.IfElse(sortByDate > 0, "ASC", "DESC")
	}
	var rows = common.AssertResultError(me.pool.Query(context.Background(), querySql))
	defer rows.Close()
	for rows.Next() {
		var row db_objects.GoalPostRow
		row.Scan(rows)
		common.AssertError(rows.Err())
		if !callback(&row) {
			break
		}
	}
}

func (me *database) getGoals() (results []db_objects.GoalRow) {
	var fields = strings.Join(getFieldNames[db_objects.GoalRow](), ",")
	var rows = common.AssertResultError(me.pool.Query(context.Background(), "SELECT "+fields+" FROM goals ORDER BY id"))
	defer rows.Close()
	for rows.Next() {
		var record db_objects.GoalRow
		record.Scan(rows)
		results = append(results, record)
	}
	return
}

func (me *database) getGoal(goalId int64) (result *db_objects.GoalRow) {
	var fields = strings.Join(getFieldNames[db_objects.GoalRow](), ",")
	var queryText = "SELECT " + fields + " FROM goals WHERE id = $1"
	var rows = common.AssertResultError(me.pool.Query(context.Background(), queryText, goalId))
	defer rows.Close()
	if rows.Next() {
		result = new(db_objects.GoalRow)
		result.Scan(rows)
	}
	return
}

func (me *database) getGoalImage(goalId int64) (imageData []byte, imageContentType string) {
	var queryText = "SELECT imageData, imageContentType FROM goals WHERE id = $1"
	var rows = common.AssertResultError(me.pool.Query(context.Background(), queryText, goalId))
	defer rows.Close()
	if rows.Next() {
		common.AssertError(rows.Scan(&imageData, &imageContentType))
	}
	return
}

func (me *database) getGoalPost(goalId int64, dateTime time.Time) (result *db_objects.GoalPostRow) {
	var queryText = "SELECT * FROM goalPosts WHERE goalId = $1 AND dateTime = $2"
	var rows = common.AssertResultError(me.pool.Query(context.Background(), queryText, goalId, dateTime.UTC().Unix()))
	defer rows.Close()
	if rows.Next() {
		result = new(db_objects.GoalPostRow)
		result.Scan(rows)
	}
	return
}

func (me *database) getGoalPostImage(goalId int64, dateTime time.Time, index int) (result *db_objects.GoalPostImageRow) {
	var queryText = "SELECT contentType, file FROM goalPostImages" +
		" WHERE goalId = $1 AND parentDateTime = $2 AND sequenceIndex = $3"
	var rows = common.AssertResultError(me.pool.Query(context.Background(), queryText, goalId, dateTime.UTC().Unix(), index))
	defer rows.Close()
	if rows.Next() {
		result = new(db_objects.GoalPostImageRow)
		common.AssertError(rows.Scan(&result.ContentType, &result.File))
	}
	return
}

func (me *database) getGoalPostImageCount(goalId int64, dateTime time.Time) (count int) {
	var queryText = "SELECT COUNT(*) FROM goalPostImages WHERE goalId = $1 AND parentDateTime = $2"
	var row = me.pool.QueryRow(context.Background(), queryText, goalId, dateTime.UTC().Unix())
	common.AssertError(row.Scan(&count))
	return
}

func (me *database) getGoalPosts(goalId int64, includePrivate bool, language language.Tag) (results []rest_objects.GoalPostHeader) {
	var titleField = "title" + db_objects.GetLanguagePostfix(language)
	var queryText = "SELECT goalId, dateTime, isPublic, type, " + titleField + " FROM goalPosts WHERE goalId = $1"
	if !includePrivate {
		queryText += " AND isPublic = TRUE"
	}
	queryText += " ORDER BY dateTime DESC"
	var rows = common.AssertResultError(me.pool.Query(context.Background(), queryText, goalId))
	defer rows.Close()
	for rows.Next() {
		var record rest_objects.GoalPostHeader
		common.AssertError(rows.Scan(&record.GoalId, &record.DateTime, &record.IsPublic, &record.Type, &record.Title))
		results = append(results, record)
	}
	return
}

func (me *database) searchGoalPosts(
	queryText string, supportedLanguage language.Tag, includePrivate bool, limit int,
) (results []*db_objects.GoalPostRow) {
	queryText = strings.ToUpper(queryText)
	queryText = common.NormalizeString(queryText)
	me.forEachGoalPost(func(row *db_objects.GoalPostRow) bool {
		if limit <= len(results) {
			return false
		}
		var isVisible = includePrivate || row.IsPublic
		if !isVisible {
			return true
		}
		var title = strings.ToUpper(row.GetTranslatedTitle(supportedLanguage))
		var text = stripHtml(strings.ToUpper(row.GetTranslatedText(supportedLanguage)))
		if strings.Contains(title, queryText) || strings.Contains(text, queryText) {
			results = append(results, row)
		}
		return true
	}, (db_objects.GoalPostRow{}).GetSelectorForLanguage(supportedLanguage), -1)
	return
}

func (me *database) migrate() {
	var query = `
		UPDATE goalPosts SET textEnglish = NULL, titleEnglish = NULL WHERE textEnglish = '';
		UPDATE goalPosts SET textGerman = NULL, titleGerman = NULL WHERE textGerman = '';
	`
	common.AssertResultError(me.pool.Exec(context.Background(), query))
}
