package server

import (
	"context"
	"time"

	_ "embed"

	"github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.postgre.sql
var dbSchemaPostgre string

type database struct {
	pool *pgxpool.Pool
}

const databaseAcquireTimeout = 5 * time.Minute

var _ pgxpool.AcquireTracer = &database{}

func (me *database) init() {
	var config = assertResultError(pgxpool.ParseConfig(requireEnvVar("POSTGRES_URL")))
	config.MaxConns = getInt32FromString(readEnvVar("POSTGRES_MAX_CONNS", "1"))
	config.ConnConfig.Tracer = me
	me.pool = assertResultError(pgxpool.NewWithConfig(context.Background(), config))
	assertResultError(me.pool.Exec(context.Background(), dbSchemaPostgre))
}

func (me *database) close() {
	if me.pool != nil {
		me.pool.Close()
		me.pool = nil
	}
}

func (me *database) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	return ctx
}

func (me *database) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func (me *database) TraceAcquireStart(ctx context.Context, pool *pgxpool.Pool, data pgxpool.TraceAcquireStartData) context.Context {
	return ctx
}

func (me *database) TraceAcquireEnd(ctx context.Context, pool *pgxpool.Pool, data pgxpool.TraceAcquireEndData) {
}
