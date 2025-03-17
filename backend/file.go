package main

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hinst/hinst-website/file_mode"
)

func readJsonFile[T any](filePath string, receiver T) T {
	var file = assertResultError(os.Open(filePath))
	defer file.Close()
	json.NewDecoder(file).Decode(receiver)
	return receiver
}

func readTextFile(filePath string) string {
	var fileContent = assertResultError(os.ReadFile(filePath))
	return string(fileContent)
}

func writeJsonFile[T any](filePath string, data T) {
	var jsonBytes = assertResultError(json.Marshal(data))
	assertError(os.WriteFile(filePath, jsonBytes, file_mode.OS_USER_RW))
}

func writeTextFile(filePath string, text string) {
	assertError(os.WriteFile(filePath, []byte(text), file_mode.OS_USER_RW))
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

func getFileNameWithoutExtension(filePath string) string {
	var fileName = filepath.Base(filePath)
	var extension = filepath.Ext(fileName)
	return strings.TrimSuffix(fileName, extension)
}

func getFilePathWithoutExtension(filePath string) string {
	var extension = filepath.Ext(filePath)
	return strings.TrimSuffix(filePath, extension)
}

func checkFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
