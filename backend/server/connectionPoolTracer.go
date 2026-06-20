package server

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muir/gwrap"
)

type ConnectionPoolTracer struct {
	timeout   time.Duration
	cancelMap gwrap.SyncMap[context.Context, context.CancelFunc]
	stackMap  gwrap.SyncMap[context.Context, string]
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
	me.stackMap.Store(ctx, string(debug.Stack()))
	return ctx
}

func (me *ConnectionPoolTracer) TraceAcquireEnd(ctx context.Context, pool *pgxpool.Pool, data pgxpool.TraceAcquireEndData) {
	var cancel, cancelExists = me.cancelMap.Load(ctx)
	if cancelExists {
		cancel()
		me.cancelMap.Delete(ctx)
		me.stackMap.Delete(ctx)
	} else {
		// should never happen
		panic("Cannot find cancel for context")
	}
}

func (me *ConnectionPoolTracer) Close() {
	var count = 0
	me.stackMap.Range(func(key context.Context, value string) bool {
		log.Println("Ongoing connection:", value)
		count += 1
		return true
	})
	if count > 0 {
		log.Println("Ongoing connection count:", count)
	}
}
