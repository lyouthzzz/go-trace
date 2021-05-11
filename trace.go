package gotrace

import (
	"context"

	"github.com/lyouthzzz/go-trace/reporter/kafka"
)

type traceKeyType int

const (
	activeSpanKey traceKeyType = iota
	remoteContextKey
)

var globalTracer Tracer

var _ Tracer = (*tracer)(nil)

type Tracer interface {
	StartSpan(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span)
}

type Config struct {
	Reporter *kafka.Reporter
}

type tracer struct {
	Name   string
	Config *Config
}

func (t *tracer) StartSpan(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span) {
	var s Span
	if parentSpan := SpanFromContext(ctx); parentSpan != nil {
		s = parentSpan.Child(name)
	} else {
		opts = append(opts, SpanTracerOption(t))
		opts = append(opts, SpanContextOption(RemoteSpanContextFromContext(ctx)))
		s = NewSpan(name, opts...)
	}
	ctx = context.WithValue(ctx, activeSpanKey, s)
	return ctx, s
}

func ContextWithRemoteSpanContext(parent context.Context, remote SpanContext) context.Context {
	return context.WithValue(parent, remoteContextKey, remote)
}

// RemoteSpanContextFromContext returns the remote span context from ctx.
func RemoteSpanContextFromContext(ctx context.Context) SpanContext {
	if sc, ok := ctx.Value(remoteContextKey).(SpanContext); ok {
		return sc
	}
	return SpanContext{}
}

func SpanFromContext(ctx context.Context) Span {
	if span, ok := ctx.Value(activeSpanKey).(Span); ok {
		return span
	}
	return nil
}

func SpanContextFromContext(ctx context.Context) SpanContext {
	if span := SpanFromContext(ctx); span != nil {
		return span.GetSpanContext()
	}
	return SpanContext{}
}

func StartSpan(ctx context.Context, name string) (context.Context, Span) {
	return StratSpanWithTracer(ctx, globalTracer, name)
}

func StratSpanWithTracer(ctx context.Context, tracer Tracer, name string) (context.Context, Span) {
	return globalTracer.StartSpan(ctx, name)
}

func SetGlobalTracer(tracer Tracer) {
	globalTracer = tracer
}

func GetGlobalTracer() Tracer {
	return globalTracer
}
