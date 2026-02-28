package server

func AssertError(err error) {
	if err != nil {
		panic(err)
	}
}

func AssertResultError[T any](result T, err error) T {
	AssertError(err)
	return result
}

func AssertCondition[T any](condition bool, exception func() T) {
	if !condition {
		panic(exception())
	}
}

func IgnoreError(err error) {
}
