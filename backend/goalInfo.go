package main

type goalInfo struct {
	title        string
	englishTitle string
	coverImage   string
}

func (goalInfo) findByTitle(array []goalInfo, title string) *goalInfo {
	for _, info := range array {
		if info.title == title {
			return &info
		}
	}
	return nil
}
