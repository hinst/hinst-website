package db_objects

import "github.com/jackc/pgx/v5"

type DbObject interface {
	GetTableName() string
	SaveToDirectory(directory string)
	GetAllColumns() []string
	Scan(rows pgx.Rows)
}

var RegisteredDbObjects []DbObject

func registerDbObject(object DbObject) int {
	RegisteredDbObjects = append(RegisteredDbObjects, object)
	return len(RegisteredDbObjects)
}
