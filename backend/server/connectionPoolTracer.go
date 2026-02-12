package server

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPoolTracer struct {
	timeout   time.Duration
	cancelMap sync.Map
}

var _ pgxpool.AcquireTracer = &ConnectionPoolTracer{}
var _ pgx.QueryTracer = &ConnectionPoolTracer{}

func (me *ConnectionPoolTracer) init() *ConnectionPoolTracer {
	return me
}

func (me *ConnectionPoolTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	return ctx
}

func (me *ConnectionPoolTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func (me *ConnectionPoolTracer) TraceAcquireStart(ctx context.Context, pool *pgxpool.Pool, data pgxpool.TraceAcquireStartData) context.Context {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, me.timeout)
	me.cancelMap.Store(ctx, cancel)
	return ctx
}

func (me *ConnectionPoolTracer) TraceAcquireEnd(ctx context.Context, pool *pgxpool.Pool, data pgxpool.TraceAcquireEndData) {
	var cancelObject, cancelExists = me.cancelMap.Load(ctx)
	if cancelExists {
		var cancel = cancelObject.(context.CancelFunc)
		cancel()
		me.cancelMap.Delete(ctx)
	} else {
		// should never happen
		panic("Cannot find cancel for context")
	}
}
