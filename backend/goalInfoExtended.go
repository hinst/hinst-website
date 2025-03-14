package main

type goalHeaderExtended struct {
	goalHeader
	LastPostDate string `json:"lastPostDate"`
}

type smartPostExtended struct {
	smartPost
	IsAutoTranslated bool `json:"isAutoTranslated"`
}
