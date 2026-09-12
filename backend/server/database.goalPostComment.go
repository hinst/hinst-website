package server

import (
	"context"
	"strings"

	"github.com/hinst/go-gophers"
	"github.com/hinst/hinst-website/server/db_objects"
)

func (me *database) saveGoalPostComment(row db_objects.GoalPostCommentRow) {
	var columnNames = row.GetAllColumns()
	var query = "INSERT INTO " + row.GetTableName() +
		" (" + strings.Join(columnNames, ",") + ")" +
		" VALUES (" + me.buildPlaceholders(len(columnNames)) + ")" +
		" ON CONFLICT (goalId, parentDateTime, dateTime, smartProgressUserId)" +
		" DO UPDATE SET username = excluded.username, text = excluded.text"
	gophers.AssertResultError(me.pool.Exec(context.Background(), query, db_objects.GetAllColumnValues(&row)...))
}
