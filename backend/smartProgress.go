package main

import "time"

const smartProgressDateTimeFormat = "2006-01-02 15:04:05"

func parseSmartProgressDateTime(text string) (time.Time, error) {
	return time.Parse(smartProgressDateTimeFormat, text)
}
