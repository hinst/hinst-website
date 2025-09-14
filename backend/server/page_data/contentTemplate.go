package page_data

import "html/template"

type ContentTemplate struct {
	Base
	Title   string
	Header  template.HTML
	Content template.HTML
}
