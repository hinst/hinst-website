package db_objects

import (
	"reflect"

	"github.com/jackc/pgx/v5"
)

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

func GetAllColumnValues(object DbObject) (items []any) {
	var columns = object.GetAllColumns()
	var value = reflect.ValueOf(object)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	for _, column := range columns {
		items = append(items, value.FieldByName(column).Interface())
	}
	return
}
