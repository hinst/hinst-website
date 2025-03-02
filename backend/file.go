package main

import (
	"encoding/json"
	"os"
)

func readJsonFile[T any](filePath string, receiver T) T {
	var fileContent = assertResultError(os.ReadFile(filePath))
	assertError(json.Unmarshal(fileContent, receiver))
	return receiver
}
