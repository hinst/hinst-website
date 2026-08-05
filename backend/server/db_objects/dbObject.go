package db_objects

import "github.com/jackc/pgx/v5"

type DbObject interface {
	GetTableName() string
	SaveToDirectory(directory string)
	GetAllColumns() []string
	Scan(rows pgx.Rows)
}

type DbObjectConstructor = func() DbObject

var DbObjects []DbObjectConstructor

func registerDbObject(constructor DbObjectConstructor) int {
	DbObjects = append(DbObjects, constructor)
	return len(DbObjects)
}
