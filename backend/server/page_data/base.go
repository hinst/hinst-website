package page_data

import "html/template"

type Base struct {
	Id         int64
	WebPath    string
	StaticPath string

	SettingsSvg template.HTML
	MenuSvg     template.HTML
	InfoSvg     template.HTML
}
