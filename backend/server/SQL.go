package server

func convertInt64ArrayToSequelString(array []int64) (text string) {
	for i, item := range array {
		if i != 0 {
			text += ", "
		}
		text += getStringFromInt64(item)
	}
	return
}
