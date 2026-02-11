package server

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPoolTracer struct {
	timeout   time.Duration
	cancelMap map[context.Context]context.CancelFunc
}

var _ pgxpool.AcquireTracer = &ConnectionPoolTracer{}
var _ pgx.QueryTracer = &ConnectionPoolTracer{}

func (me *ConnectionPoolTracer) init() *ConnectionPoolTracer {
	me.cancelMap = make(map[context.Context]context.CancelFunc)
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
	me.cancelMap[ctx] = cancel
	return ctx
}

func (me *ConnectionPoolTracer) TraceAcquireEnd(ctx context.Context, pool *pgxpool.Pool, data pgxpool.TraceAcquireEndData) {
	var cancel, cancelExists = me.cancelMap[ctx]
	if cancelExists {
		cancel()
		delete(me.cancelMap, ctx)
	} else {
		// should never happen
		panic("Cannot find cancel for context")
	}
}
