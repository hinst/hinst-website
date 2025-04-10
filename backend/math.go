package main

func absInt64(number int64) int64 {
	if number < 0 {
		return -number
	}
	return number
}
