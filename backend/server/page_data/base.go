package page_data

import "html/template"

type Base struct {
	Id         int64
	WebPath    string
	StaticPath string
	// For server side rendering, should be empty. For static files, should be .jpg
	JpegExtension string
	// For server side rendering, should be empty. For static files, should be .html
	HtmlExtension string

	SettingsSvg template.HTML
	MenuSvg     template.HTML
}
