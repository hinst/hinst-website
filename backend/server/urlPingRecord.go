package server

type urlPingRecord struct {
	url     string
	service int
	doneAt  int64 // Unix seconds UTC
}
