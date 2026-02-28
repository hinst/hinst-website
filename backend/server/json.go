package server

import (
	"encoding/json"

	"github.com/hinst/go-common"
)

func encodeJson[T any](value T) []byte {
	return common.AssertResultError(json.Marshal(value))
}

func decodeJson[T any](data []byte, value T) T {
	common.AssertError(json.Unmarshal(data, value))
	return value
}
