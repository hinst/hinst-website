package server

import (
	"database/sql"
	"time"

	"github.com/hinst/go-common"
)

type riddleRow struct {
	id        int
	product   int
	createdAt time.Time
}

func (me *riddleRow) scan(row *sql.Rows) {
	var createdAt int64
	common.AssertError(row.Scan(&me.id, &me.product, &createdAt))
	me.createdAt = time.Unix(createdAt, 0)
}
