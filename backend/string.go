package main

import "strconv"

func getIntFromString(text string) int {
	return assertResultError(strconv.Atoi(text))
}
