package page_data

import "html/template"

type GoalCard struct {
	Base
	Id    int64
	Title string
	Image template.URL
}
