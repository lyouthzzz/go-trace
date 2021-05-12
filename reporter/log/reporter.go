package log

import (
	gotrace "github.com/lyouthzzz/go-trace"
	"go.uber.org/zap"
)

var _ gotrace.Reporter = (*Reporter)(nil)

type Reporter struct {
	logger *zap.Logger
}

func NewReporter(logger *zap.Logger) gotrace.Reporter {
	return &Reporter{logger: logger}
}

func (r *Reporter) Send(span gotrace.Span) {
	var (
		spanModel *gotrace.SpanModel
		ok        bool
	)
	if spanModel, ok = span.(*gotrace.SpanModel); !ok {
		return
	}
	spanCtx := spanModel.GetSpanContext()

	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("globalTicket", spanCtx.GetGlobalTicketId()),
		zap.String("monitorTrackId", spanCtx.GetMonitorId()),
		zap.String("parentRpcId", spanCtx.GetParentRpcId()),
		zap.String("rpcId", spanCtx.GetRpcId()),
		zap.Time("invokeTime", spanModel.StartTime),
		// zap.String("invokeTime", spanModel.StartTime.Format("2006-01-02 15:04:05")),
		zap.Int64("elapsed", int64(spanModel.Duration)),
		zap.String("invokeType", string(spanModel.Kind)),
		zap.Int("isSuccess", int(spanModel.Status)),
	)

	for k, v := range spanModel.Attributes {
		fields = append(fields, zap.Any(k, v))
	}
	r.logger.Info(spanModel.Name, fields...)
}

func (r *Reporter) Close() error {
	return nil
}
