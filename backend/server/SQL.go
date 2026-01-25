package server

import "strings"

func convertInt64ArrayToSequelString(array []int64) (text string) {
	for i, item := range array {
		if i != 0 {
			text += ", "
		}
		text += getStringFromInt64(item)
	}
	return
}

func escapeLikeString(text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, "%", "\\%")
	text = strings.ReplaceAll(text, "'", "''")
	return text
}
