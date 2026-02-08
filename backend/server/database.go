package server

import (
	"context"
	"database/sql"
	"time"

	_ "embed"

	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

// Measured in milliseconds
var dbTimeout = time.Hour / time.Millisecond

//go:embed schema.postgre.sql
var dbSchemaPostgre string

type database struct {
	pool *pgxpool.Pool
}

func (me *database) init() {
	var config = assertResultError(pgxpool.ParseConfig(requireEnvVar("POSTGRES_URL")))
	me.pool = assertResultError(pgxpool.NewWithConfig(context.Background(), config))
	assertResultError(me.pool.Exec(context.Background(), dbSchemaPostgre))
}

func (me *database) close(db *sql.DB) *sql.DB {
	if db != nil {
		assertError(db.Close())
	}
	return nil
}
