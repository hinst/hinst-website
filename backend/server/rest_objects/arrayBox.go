package rest_objects

type ArrayBox[T any] struct {
	Items []T `json:"items"`
}

func NewArrayBox[T any](Items []T) ArrayBox[T] {
	return ArrayBox[T]{Items: Items}
}

func NewArrayBoxPtr[T any](Items []T) *ArrayBox[T] {
	return &ArrayBox[T]{Items: Items}
}
