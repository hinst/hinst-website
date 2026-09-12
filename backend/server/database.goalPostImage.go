package server

import (
	"context"
	"strings"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/db_objects"
)

func (me *database) saveGoalPostImage(row db_objects.GoalPostImageRow) {
	var columnNames = row.GetAllColumns()
	var query = "INSERT INTO " + row.GetTableName() +
		" (" + strings.Join(columnNames, ",") + ")" +
		" VALUES (" + me.buildPlaceholders(len(columnNames)) + ")" +
		" ON CONFLICT (goalId, parentDateTime, sequenceIndex)" +
		" DO UPDATE SET contentType = excluded.contentType, file = excluded.file"
	gophers.AssertResultError(me.pool.Exec(context.Background(), query, db_objects.GetAllColumnValues(&row)...))
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
