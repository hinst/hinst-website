package server

import (
	"context"
	"time"

	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.postgre.sql
var dbSchemaPostgre string

type database struct {
	pool *pgxpool.Pool
}

func (me *database) init() {
	var config = assertResultError(pgxpool.ParseConfig(requireEnvVar("POSTGRES_URL")))
	config.MaxConns = getInt32FromString(readEnvVar("POSTGRES_MAX_CONNS", "2"))
	config.ConnConfig.Tracer = (&ConnectionPoolTracer{timeout: 1 * time.Minute}).init()
	me.pool = assertResultError(pgxpool.NewWithConfig(context.Background(), config))
	assertResultError(me.pool.Exec(context.Background(), dbSchemaPostgre))
}

func (me *database) close() {
	if me.pool != nil {
		me.pool.Close()
		me.pool = nil
	}
}
