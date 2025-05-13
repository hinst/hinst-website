package main

import "time"

type riddleItem struct {
	id        int
	keys      []int
	product   int
	createdAt time.Time
}
