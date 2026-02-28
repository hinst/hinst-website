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
	var config = AssertResultError(pgxpool.ParseConfig(requireEnvVar("POSTGRES_URL")))
	config.MaxConns = getInt32FromString(readEnvVar("POSTGRES_MAX_CONNS", "2"))
	config.ConnConfig.Tracer = (&ConnectionPoolTracer{timeout: 1 * time.Minute}).init()
	me.pool = AssertResultError(pgxpool.NewWithConfig(context.Background(), config))
	AssertResultError(me.pool.Exec(context.Background(), dbSchemaPostgre))
}

func (me *database) close() {
	if me.pool != nil {
		me.pool.Close()
		me.pool = nil
	}
}
