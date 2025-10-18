package server

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/hinst/hinst-website/server/file_mode"
)

func readJsonFile[T any](filePath string, receiver T) T {
	var file = assertResultError(os.Open(filePath))
	defer ioCloseSilently(file)
	assertError(json.NewDecoder(file).Decode(receiver))
	return receiver
}

func readTextFile(filePath string) string {
	return string(readBytesFile(filePath))
}

func readBytesFile(filePath string) []byte {
	var fileContent = assertResultError(os.ReadFile(filePath))
	return fileContent
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

func checkFileExists(filePath string) bool {
	var info, err = os.Stat(filePath)
	return err == nil && !info.IsDir()
}

func checkDirectoryExists(directoryPath string) bool {
	var info, err = os.Stat(directoryPath)
	return err == nil && info.IsDir()
}

func copyFile(destinationPath string, sourcePath string) (size int64) {
	var sourceFile = assertResultError(os.Open(sourcePath))
	defer ioCloseSilently(sourceFile)
	var destinationFile = assertResultError(os.Create(destinationPath))
	defer ioClose(destinationFile)
	return assertResultError(io.Copy(destinationFile, sourceFile))
}

// https://stackoverflow.com/questions/29505089/how-can-i-compare-two-files-in-golang
func checkFilesEqual(file1, file2 string) bool {
	const chunkSize = 64000
	f1, err := os.Open(file1)
	if err != nil {
		return false
	}
	defer ioCloseSilently(f1)

	f2, err := os.Open(file2)
	if err != nil {
		return false
	}
	defer ioCloseSilently(f2)

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				panic("File comparison error")
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func checkTextFilesEqual(file1, file2 string) bool {
	var text1 = ""
	if checkFileExists(file1) {
		text1 = readTextFile(file1)
	}
	var text2 = ""
	if checkFileExists(file2) {
		text2 = readTextFile(file2)
	}
	// ignore line breaks
	text1 = strings.ReplaceAll(text1, "\r\n", "\n")
	text2 = strings.ReplaceAll(text2, "\r\n", "\n")
	// trim whitespace
	text1 = strings.TrimSpace(text1)
	text2 = strings.TrimSpace(text2)
	return text1 == text2
}
