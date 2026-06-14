package server

import (
	"context"
	"log"
	"time"

	_ "embed"

	"github.com/hinst/go-gophers"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.postgre.sql
var dbSchemaPostgre string

type database struct {
	pool *pgxpool.Pool
}

func (me *database) init() {
	var config = gophers.AssertResultError(pgxpool.ParseConfig(gophers.RequireEnvVar("POSTGRES_URL")))
	config.MaxConns = gophers.GetInt32FromString(gophers.ReadEnvVar("POSTGRES_MAX_CONNS", "2"))
	config.ConnConfig.Tracer = (&ConnectionPoolTracer{timeout: 1 * time.Minute}).init()
	me.pool = gophers.AssertResultError(pgxpool.NewWithConfig(context.Background(), config))
	gophers.AssertResultError(me.pool.Exec(context.Background(), dbSchemaPostgre))
	gophers.InstallShutdownReceiver(func() {
		me.close()
	})
}

func (me *database) close() {
	log.Print("Closing...")
	if me.pool != nil {
		me.pool.Close()
		me.pool = nil
	}
}
