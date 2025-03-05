package main

import (
	"encoding/json"
	"io/fs"
	"os"
	"slices"
)

func readJsonFile[T any](filePath string, receiver T) T {
	var fileContent = assertResultError(os.ReadFile(filePath))
	assertError(json.Unmarshal(fileContent, receiver))
	return receiver
}

func writeJsonFile[T any](filePath string, data T) {
	var jsonBytes = assertResultError(json.Marshal(data))
	assertError(os.WriteFile(filePath, jsonBytes, 0644))
}

func sortFilesByName(files []fs.DirEntry) {
	slices.SortFunc(files, func(a fs.DirEntry, b fs.DirEntry) int {
		if a.Name() < b.Name() {
			return -1
		} else if a.Name() == b.Name() {
			return 0
		} else {
			return 1
		}
	})
}
