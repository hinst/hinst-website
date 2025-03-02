package main

import "encoding/json"

func encodeJson[T any](value T) []byte {
	return assertResultError(json.Marshal(value))
}
