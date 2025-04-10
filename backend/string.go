package main

import "strconv"

func getIntFromString(text string) int {
	return assertResultError(strconv.Atoi(text))
}

func getStringFromInt64(number int64) string {
	return strconv.FormatInt(number, 10)
}

func getStringFromBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}
