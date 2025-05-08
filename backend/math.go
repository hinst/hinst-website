package main

func absInt64(number int64) int64 {
	if number < 0 {
		return -number
	}
	return number
}

func multiplyLimited(a int, b int, limit int) int {
	var result = (int64(a) * int64(b)) % int64(limit)
	return int(result)
}
