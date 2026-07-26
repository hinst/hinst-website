package page_data

import "html/template"

type Base struct {
	Id         int64
	WebPath    string // Full site URL prefix: "" or "/hinst-website" or "/blog"
	LangPath   string // Per-language segment appended after WebPath: "", "/de", "/ru"

	SettingsSvg template.HTML
	MenuSvg     template.HTML
	InfoSvg     template.HTML
}
