package server

import "encoding/json"

func encodeJson[T any](value T) []byte {
	return AssertResultError(json.Marshal(value))
}

func decodeJson[T any](data []byte, value T) T {
	AssertError(json.Unmarshal(data, value))
	return value
}
