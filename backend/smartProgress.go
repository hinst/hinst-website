package main

import "time"

const smartProgressTimeFormat = "2006-01-02 15:04:05"
const storedGoalFileTimeFormat = "2006-01-02_15-04-05"

func parseSmartProgressDateTime(text string) (time.Time, error) {
	return time.Parse(smartProgressTimeFormat, text)
}
