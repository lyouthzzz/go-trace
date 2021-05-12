package gotrace

import (
	"context"

	"github.com/google/uuid"
	"github.com/lyouthzzz/go-trace/snowflake"
)

const (
	GloblaTicketIdKey = "X-Global-Ticket-ID"
	MonitorIdKey      = "X-Monitor-ID"
	RpcIdKey          = "X-RPC-ID"
	ChildRpcIdKey     = "X-Child-RPC-ID"
)

type Progagator interface {
	Extract(ctx context.Context, carrier Carrier) context.Context
	Inject(ctx context.Context, carrier Carrier)
}

type propagator struct{}

var (
	_ Progagator = (*propagator)(nil)

	defaultPropagator = propagator{}
)

func (h propagator) Extract(ctx context.Context, carrier Carrier) context.Context {
	gid := carrier.Get(GloblaTicketIdKey)
	rid := carrier.Get(RpcIdKey)
	cid := carrier.Get(ChildRpcIdKey)
	if gid == "" {
		gid = snowflake.NewID().String()
	}
	if rid == "" {
		rid = "0.1"
	}
	if cid == "" {
		cid = rid + ".1"
	}

	sc := SpanContext{
		globalTicketId: gid,
		parentRpcId:    rid,
		rpcId:          cid,
		monitorId:      uuid.New().String(),
	}
	ctx = ContextWithRemoteSpanContext(ctx, sc)
	return ctx
}

func (h propagator) Inject(ctx context.Context, carrier Carrier) {
	sc := SpanContextFromContext(ctx)

	carrier.Set(GloblaTicketIdKey, sc.globalTicketId)
	carrier.Set(RpcIdKey, sc.parentRpcId)
	carrier.Set(ChildRpcIdKey, sc.rpcId)
}

func GetPropagator() Progagator {
	return defaultPropagator
}
