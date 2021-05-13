package middleware

import (
	"context"

	gotrace "github.com/lyouthzzz/go-trace"
	"github.com/lyouthzzz/go-trace/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientTracingInterceptor() grpc.UnaryClientInterceptor {
	propagator := gotrace.GetPropagator()
	tracer := gotrace.GetGlobalTracer()

	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.Pairs()
		}
		carrier := gotrace.GRPCCarrier(md)
		ctx = metadata.NewOutgoingContext(ctx, md)

		ctx, span := tracer.StartSpan(ctx, method,
			gotrace.SpanKindOption(gotrace.SpanKindConsumer),
			gotrace.SpanTracerOption(tracer))
		span.SetAttributes(attribute.RpcType("GRPC"))
		span.SetAttributes(attribute.KV("method", method))
		defer span.End()

		propagator.Inject(ctx, carrier)

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			span.SetStatus(gotrace.SpanStatusFail)
			span.AddError(err)
		}
		return err
	}
}

func ServerUnaryTracingInterceptor() grpc.UnaryServerInterceptor {
	propagator := gotrace.GetPropagator()
	tracer := gotrace.GetGlobalTracer()

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}
		carrier := gotrace.GRPCCarrier(md)
		if err != nil {
			return handler(ctx, req)
		}
		ctx = propagator.Extract(ctx, carrier)

		ctx, span := tracer.StartSpan(ctx, info.FullMethod,
			gotrace.SpanKindOption(gotrace.SpanKindProvider),
			gotrace.SpanTracerOption(tracer))
		span.SetAttributes(attribute.RpcType("GRPC"))
		defer span.End()

		reply, err := handler(ctx, req)
		if err != nil {
			span.SetStatus(gotrace.SpanStatusFail)
			span.AddError(err)
		}
		return reply, err
	}
}
