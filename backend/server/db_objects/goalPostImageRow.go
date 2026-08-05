package db_objects

type GoalPostImageRow struct {
	ContentType string
	File        []byte
}

var _ = registerDbObject(GoalPostImageRow{})

func (GoalPostImageRow) GetTableName() string {
	return "goalPostImages"
}

func (GoalPostImageRow) SaveToDirectory(directory string) {
	//TODO
}
