package main

import (
	"os"
	"regexp"
	"sort"
	"time"
)

const smartProgressTimeFormat = "2006-01-02 15:04:05"
const storedGoalFileTimeFormat = "2006-01-02_15-04-05"

var goalIdStringMatcher = regexp.MustCompile(`^\d{1,10}$`)
var goalFileNameMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d`)

const savedGoalHeaderFileName = "_header.json"

func parseSmartProgressDate(text string) (time.Time, error) {
	return time.Parse(smartProgressTimeFormat, text)
}

func parseStoredGoalFileDate(text string) (time.Time, error) {
	return time.Parse(storedGoalFileTimeFormat, text)
}

func getGoalFiles(goalDirectory string) (files []string) {
	var filesRaw = assertResultError(os.ReadDir(goalDirectory))
	for _, file := range filesRaw {
		if !file.IsDir() && goalFileNameMatcher.MatchString(file.Name()) {
			files = append(files, file.Name())
		}
	}
	sort.Strings(files)
	return
}
