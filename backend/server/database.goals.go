package server

import (
	"context"
	"strconv"
	"strings"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/db_objects"
)

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
