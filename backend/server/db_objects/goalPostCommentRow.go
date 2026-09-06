package db_objects

import (
	"os"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/go-gophers/file_mode"
	"github.com/hinst/hinst-website/server/base"
	"github.com/jackc/pgx/v5"
)

type GoalPostCommentRow struct {
	GoalId int64
	/* Unix seconds UTC */
	ParentDateTime int64
	/* Unix seconds UTC */
	DateTime            int64
	SmartProgressUserId int64
	Username            string
	Text                string
}

var _ = registerDbObject(func() DbObject { return new(GoalPostCommentRow) })

func (GoalPostCommentRow) GetTableName() string {
	return "goalPostComments"
}

func (GoalPostCommentRow) GetAllColumns() []string {
	return gophers.GetFieldNames[GoalPostCommentRow]()
}

func (me *GoalPostCommentRow) Scan(rows pgx.Rows) {
	gophers.AssertError(rows.Scan(
		&me.GoalId,
		&me.ParentDateTime,
		&me.DateTime,
		&me.SmartProgressUserId,
		&me.Username,
		&me.Text,
	))
}

func (me GoalPostCommentRow) SaveToDirectory(directory string) {
	directory += "/" + gophers.GetStringFromInt64(me.GoalId) + "/" +
		me.GetParentDateTime().UTC().Format("2006-01-02_15-04-05")
	gophers.AssertError(os.MkdirAll(directory, file_mode.USER_RWX))
	var filePath = directory + "/" + me.GetDateTime().UTC().Format("2006-01-02_15-04-05") + ".yaml"
	gophers.WriteBytesFile(filePath, base.EncodeYaml(me))
}

func (me *GoalPostCommentRow) GetParentDateTime() (result time.Time) {
	return time.Unix(me.ParentDateTime, 0)
}

func (me *GoalPostCommentRow) GetDateTime() (result time.Time) {
	return time.Unix(me.DateTime, 0)
}
