package server

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hinst/hinst-website/server/file_mode"
)

func readJsonFile[T any](filePath string, receiver T) T {
	var file = assertResultError(os.Open(filePath))
	defer file.Close()
	json.NewDecoder(file).Decode(receiver)
	return receiver
}

func readJsonFiles[T any](filePaths []string, threadCount int) (items []*T) {
	type filePathItem struct {
		filePath string
		index    int
	}
	var inputs = make(chan filePathItem)
	var ready = make(chan struct{}, threadCount)
	items = make([]*T, len(filePaths))
	var reader = func() {
		for input := range inputs {
			var item T
			readJsonFile(input.filePath, &item)
			items[input.index] = &item
		}
		ready <- struct{}{}
	}
	for range threadCount {
		go reader()
	}
	for index, filePath := range filePaths {
		inputs <- filePathItem{filePath, index}
	}
	close(inputs)
	for range threadCount {
		<-ready
	}
	close(ready)
	return
}

func readTextFile(filePath string) string {
	var fileContent = assertResultError(os.ReadFile(filePath))
	return string(fileContent)
}

func writeBytesFile(filePath string, data []byte) {
	assertError(os.WriteFile(filePath, data, file_mode.OS_USER_RW))
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

func checkFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
