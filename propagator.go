package gotrace

import (
	"context"
)

const (
	globlaTicketIdKey = "X-Global-Ticket-ID"
	monitorIdKey      = "X-Monitor-ID"
	rpcIdKey          = "X-RPC-ID"
	childRpcIdKey     = "X-Child-RPC-ID"
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
	gid := carrier.Get(globlaTicketIdKey)
	rid := carrier.Get(rpcIdKey)
	cid := carrier.Get(childRpcIdKey)

	sc := SpanContext{
		globalTicketId: gid,
		parentRpcId:    rid,
		rpcId:          cid,
		monitorId:      "monitorUUID",
	}
	ctx = ContextWithRemoteSpanContext(ctx, sc)
	return ctx
}

func (h propagator) Inject(ctx context.Context, carrier Carrier) {
	sc := SpanContextFromContext(ctx)

	carrier.Set(globlaTicketIdKey, sc.globalTicketId)
	carrier.Set(rpcIdKey, sc.parentRpcId)
	carrier.Set(childRpcIdKey, sc.rpcId)
}

func GetPropagator() Progagator {
	return defaultPropagator
}
