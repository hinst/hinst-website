package main

import (
	"database/sql"
	"time"
)

type riddleItem struct {
	Id        int `json:"id"`
	keys      []int
	Product   int `json:"product"`
	createdAt time.Time
}

func (me *riddleItem) scan(row *sql.Rows) {
	var createdAt int64
	assertError(row.Scan(&me.Id, &me.Product, &createdAt))
	me.createdAt = time.Unix(createdAt, 0)
}
