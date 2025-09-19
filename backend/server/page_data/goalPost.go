package page_data

import "html/template"

type GoalPost struct {
	Base
	GoalId   int64
	DateTime int64
	Text     template.HTML
	Images   []int
}
