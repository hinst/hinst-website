package db_objects

import (
	"github.com/hinst/go-gophers"
)

type GoalPostImageRow struct {
	ContentType string
	File        []byte
}

var _ = registerDbObject(GoalPostImageRow{})

func (GoalPostImageRow) GetTableName() string {
	return "goalPostImages"
}

func (GoalPostImageRow) GetAllColumns() []string {
	return gophers.GetFieldNames[GoalPostImageRow]()
}

func (GoalPostImageRow) SaveToDirectory(directory string) {
	//TODO
}
