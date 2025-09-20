package page_data

import "html/template"

type Base struct {
	Id          int64
	WebPath     string
	ApiPath     string
	SettingsSvg template.HTML
}
