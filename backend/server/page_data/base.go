package page_data

import "html/template"

type Base struct {
	Id          int64
	WebPath     string
	SettingsSvg template.HTML
}
