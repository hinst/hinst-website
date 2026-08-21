package rest_objects

// Response container for Huma Rest framework
type Response[T any] struct {
	Body         T
	ContentType  string `header:"Content-Type"`
	CacheControl string `header:"Cache-Control"`
}

func NewSimpleResponse[T any](Item T) *Response[T] {
	return &Response[T]{Body: Item}
}
