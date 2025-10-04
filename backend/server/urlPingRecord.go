package server

type urlPingRecord struct {
	Url     string
	Service int
	DoneAt  int64 // Unix seconds UTC
}
