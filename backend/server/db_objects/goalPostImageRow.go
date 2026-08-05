package db_objects

func (GoalPostImageRow) GetTableName() string {
	return "goalPostImages"
}

type GoalPostImageRow struct {
	ContentType string
	File        []byte
}
