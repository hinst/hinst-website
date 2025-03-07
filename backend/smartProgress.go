package main

import "time"

const SMART_PROGRESS_DATE_TIME_FORMAT = "2006-01-02 15:04:05"

func parseSmartProgressDateTime(text string) (time.Time, error) {
	return time.Parse(SMART_PROGRESS_DATE_TIME_FORMAT, text)
}
