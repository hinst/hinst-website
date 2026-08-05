package db_objects

type DbObject interface {
	GetTableName() string
	SaveToDirectory(directory string)
	GetAllColumns() []string
}

var RegisteredDbObjects []DbObject

func registerDbObject(object DbObject) int {
	RegisteredDbObjects = append(RegisteredDbObjects, object)
	return len(RegisteredDbObjects)
}
