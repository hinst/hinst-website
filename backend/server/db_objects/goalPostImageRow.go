package db_objects

import (
	"os"

	"github.com/hinst/go-gophers"
	"github.com/hinst/go-gophers/file_mode"
	"github.com/jackc/pgx/v5"
)

type GoalPostImageRow struct {
	GoalId         int64
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
		gophers.GetStringFromInt64(me.ParentDateTime)
	gophers.AssertError(os.MkdirAll(directory, file_mode.USER_RWX))
	var extension = gophers.AssertResultError(
		gophers.MimeContentType{}.GetFileExtension(me.ContentType))
	var filePath = directory + "/" + gophers.GetStringFromInt64(me.SequenceIndex) + extension
	gophers.WriteBytesFile(filePath, me.File)
}
