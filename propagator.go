package gotrace

import (
	"context"

	"github.com/google/uuid"
)

const (
	GloblaTicketIdKey = "globalTicket"
	ParentRpcId       = "parentRpcId"
	RpcEntryUrl       = "rpcEntryUrl"
	RpcIndex          = "rpcIndex"
	RpcId             = "rpcId"
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
	pid := carrier.Get(ParentRpcId)
	rid := carrier.Get(RpcId)
	ri := carrier.Get(RpcIndex)
	if gid == "" {
		gid = uuid.New().String()
	}
	if pid == "" {
		pid = "0.1"
	}
	if ri == "" {
		ri = "1"
	}
	if rid == "" {
		rid = pid + "." + ri
	}

	sc := SpanContext{
		globalTicketId: gid,
		parentRpcId:    pid,
		rpcIndex:       ri,
		rpcId:          rid,
		monitorId:      uuid.New().String(),
	}
	ctx = ContextWithRemoteSpanContext(ctx, sc)
	return ctx
}

func (h propagator) Inject(ctx context.Context, carrier Carrier) {
	sc := SpanContextFromContext(ctx)

	carrier.Set(GloblaTicketIdKey, sc.globalTicketId)
	carrier.Set(ParentRpcId, sc.parentRpcId)
	carrier.Set(RpcIndex, sc.rpcIndex)
	carrier.Set(RpcId, sc.rpcId)
}

func GetPropagator() Progagator {
	return defaultPropagator
}
