package main

type goalInfo struct {
	title        string
	englishTitle string
	coverImage   string
}

func findEnglishTitle(array []goalInfo, title string) string {
	for _, info := range array {
		if info.title == title {
			return info.englishTitle
		}
	}
	return ""
}
