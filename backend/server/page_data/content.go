package page_data

import "html/template"

type Content struct {
	Base
	Title       string
	Description string
	Header      template.HTML
	Content     template.HTML
}
