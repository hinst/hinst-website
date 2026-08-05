package db_objects

func (GoalPostImageRow) getTableName() string {
	return "goalPostImages"
}

type GoalPostImageRow struct {
	ContentType string
	File        []byte
}
