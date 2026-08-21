package rest_objects

// Response container for Huma Rest framework
type Response[T any] struct {
	Body         T
	CacheControl string `header:"Cache-Control"`
	ContentType  string `header:"Content-Type"`
}

func NewSimpleResponse[T any](Item T) *Response[T] {
	return &Response[T]{Body: Item}
}
