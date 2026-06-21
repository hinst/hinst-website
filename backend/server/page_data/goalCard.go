package page_data

import "html/template"

type GoalCard struct {
	Base
	Id    int64
	Title string
	// Base64 encoded data
	Image template.URL
}
