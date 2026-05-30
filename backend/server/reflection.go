package server

import "reflect"

func getFieldNames[T any]() []string {
	// Create a nil pointer of type *T and get its element type
	t := reflect.TypeOf((*T)(nil)).Elem()

	names := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		names[i] = t.Field(i).Name
	}

	return names
}
