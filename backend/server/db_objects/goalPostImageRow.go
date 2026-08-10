package db_objects

import (
	"os"
	"time"

	"github.com/hinst/go-gophers"
	"github.com/hinst/go-gophers/file_mode"
	"github.com/jackc/pgx/v5"
)

type GoalPostImageRow struct {
	GoalId int64
	/* Unix seconds UTC */
	ParentDateTime int64
	SequenceIndex  int64
	ContentType    string
	File           []byte
}

func (me *GoalPostImageRow) Scan(rows pgx.Rows) {
	gophers.AssertError(rows.Scan(&me.GoalId, &me.ParentDateTime, &me.SequenceIndex, &me.ContentType, &me.File))
}

var _ = registerDbObject(func() DbObject { return new(GoalPostImageRow) })

func (GoalPostImageRow) GetTableName() string {
	return "goalPostImages"
}

func (GoalPostImageRow) GetAllColumns() []string {
	return gophers.GetFieldNames[GoalPostImageRow]()
}

func (me GoalPostImageRow) SaveToDirectory(directory string) {
	directory += "/" + gophers.GetStringFromInt64(me.GoalId) + "/" +
		me.GetParentDateTime().UTC().Format("2006-01-02_15-04-05")
	gophers.AssertError(os.MkdirAll(directory, file_mode.USER_RWX))
	var extension = gophers.AssertResultError(
		gophers.MimeContentType{}.GetFileExtension(me.ContentType))
	var filePath = directory + "/" + gophers.GetStringFromInt64(me.SequenceIndex) + extension
	gophers.WriteBytesFile(filePath, me.File)
}

func (me *GoalPostImageRow) GetParentDateTime() time.Time {
	return time.Unix(me.ParentDateTime, 0)
}
