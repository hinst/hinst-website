package db_objects

import (
	"github.com/hinst/go-gophers"
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

func (GoalPostImageRow) SaveToDirectory(directory string) {
	//TODO
}
