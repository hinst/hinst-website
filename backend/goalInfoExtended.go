package main

type goalHeaderExtended struct {
	goalHeader
	LastPostDate string `json:"lastPostDate"`
}
