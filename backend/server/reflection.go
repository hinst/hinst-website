package server

import "reflect"

func getFieldNames[T any]() []string {
	var theType = reflect.TypeFor[T]()
	var names = make([]string, theType.NumField())
	for i := 0; i < theType.NumField(); i++ {
		names[i] = theType.Field(i).Name
	}
	return names
}
