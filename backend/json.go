package main

import "github.com/goccy/go-json"

func encodeJson[T any](value T) []byte {
	return assertResultError(json.Marshal(value))
}

func decodeJson[T any](data []byte, value T) T {
	assertError(json.Unmarshal(data, value))
	return value
}
