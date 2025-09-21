package page_data

import "html/template"

type Base struct {
	Id            int64
	WebPath       string
	StaticPath    string
	JpegExtension string
	SettingsSvg   template.HTML
}
