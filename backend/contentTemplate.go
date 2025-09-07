package main

import "html/template"

type ContentTemplate struct {
	BaseTemplate
	Content template.HTML
}
