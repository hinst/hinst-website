package rest_objects

type BodyBox[T any] struct {
	Body T
}

func NewBodyBox[T any](Item T) *BodyBox[T] {
	return &BodyBox[T]{Body: Item}
}
