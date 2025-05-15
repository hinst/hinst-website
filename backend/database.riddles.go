package main

import "context"

func (me *database) insertRiddle(item *riddleRow) {
	var db = me.open()
	defer me.close(db)
	var row = db.QueryRow(
		"INSERT INTO riddles (id, product, createdAt) VALUES (null, ?, ?) RETURNING id",
		item.product,
		item.createdAt.UTC().Unix(),
	)
	assertError(row.Err())
	assertError(row.Scan(&item.id))
}

func (me *database) processRiddle(id int, product int, function func(item *riddleRow)) {
	var db = me.open()
	var transaction = assertResultError(db.BeginTx(context.Background(), nil))
	defer func() {
		assertError(transaction.Commit())
		me.close(db)
	}()
	var row = assertResultError(transaction.Query(
		"SELECT id, product, createdAt FROM riddles WHERE id = ? AND product = ?",
		id, product,
	))
	var item *riddleRow
	if row.Next() {
		item = new(riddleRow)
		item.scan(row)
	}
	function(item)
	if item != nil {
		assertResultError(transaction.Exec("DELETE FROM riddles WHERE id = ?", item.id))
	}
}
