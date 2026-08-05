package db_objects

type DbObject interface {
	getTableName() string
}
