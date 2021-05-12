package main

import (
	"context"

	gotrace "github.com/lyouthzzz/go-trace"
	"github.com/lyouthzzz/go-trace/attribute"
	"github.com/lyouthzzz/go-trace/reporter/log"
)

func main() {
	tracer := gotrace.NewTracer("tracer", log.NewReporter(log.NewStdoutLogger(log.InfoLevel)))
	gotrace.SetGlobalTracer(tracer)

	_, span := gotrace.StartSpan(context.Background(), "span", gotrace.SpanKindOption(gotrace.SpanKindConsumer))

	span.SetStatus(gotrace.SpanStatusSuccess)
	span.SetAttributes(attribute.RpcType("GRPC"))
	span.SetAttributes(attribute.KV("A", "AA"))
	span.SetAttributes(attribute.KV("B", "BB"))

	span.End()
}
