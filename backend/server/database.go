package server

import (
	"context"
	"time"

	_ "embed"

	"github.com/hinst/go-gophers"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.postgre.sql
var dbSchemaPostgre string

type database struct {
	tracer *ConnectionPoolTracer
	pool   *pgxpool.Pool
}

func (me *database) init() {
	var config = gophers.AssertResultError(pgxpool.ParseConfig(gophers.RequireEnvVar("POSTGRES_URL")))
	config.MaxConns = gophers.GetInt32FromString(gophers.ReadEnvVar("POSTGRES_MAX_CONNS", "2"))
	me.tracer = (&ConnectionPoolTracer{timeout: 1 * time.Minute}).init()
	config.ConnConfig.Tracer = me.tracer
	me.pool = gophers.AssertResultError(pgxpool.NewWithConfig(context.Background(), config))
	gophers.AssertResultError(me.pool.Exec(context.Background(), dbSchemaPostgre))
}

func (me *database) close() {
	if me.tracer != nil {
		me.tracer.Close()
		me.tracer = nil
	}
	if me.pool != nil {
		me.pool.Close()
		me.pool = nil
	}
}
