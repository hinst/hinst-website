package main

import "html/template"

type ContentTemplate struct {
	BaseTemplate
	Title   string
	Header  template.HTML
	Content template.HTML
}
