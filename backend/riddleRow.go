package main

import (
	"database/sql"
	"time"
)

type riddleRow struct {
	id        int
	product   int
	createdAt time.Time
}

func (me *riddleRow) scan(row *sql.Rows) {
	var createdAt int64
	assertError(row.Scan(&me.id, &me.product, &createdAt))
	me.createdAt = time.Unix(createdAt, 0)
}
