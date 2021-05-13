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
		sm *gotrace.SpanModel
		ok bool
	)
	if sm, ok = span.(*gotrace.SpanModel); !ok {
		return
	}
	spanCtx := sm.GetSpanContext()

	var fields = make([]zap.Field, 0)
	fields = append(fields, zap.String("globalTicket", spanCtx.GetGlobalTicketId()),
		zap.String("monitorTrackId", spanCtx.GetMonitorId()),
		zap.String("parentRpcId", spanCtx.GetParentRpcId()),
		zap.String("rpcId", spanCtx.GetRpcId()),
		zap.Time("invokeTime", sm.StartTime),
		// zap.String("invokeTime", spanModel.StartTime.Format("2006-01-02 15:04:05")),
		zap.Int64("elapsed", int64(sm.Duration)),
		zap.String("invokeType", string(sm.Kind)),
		zap.Int("isSuccess", int(sm.Status)),
	)

	if sm.Errs != nil && len(sm.Errs) > 0 {
		fields = append(fields, zap.Error(sm.Errs[0]))
	}
	for k, v := range sm.Attributes {
		fields = append(fields, zap.Any(k, v))
	}
	r.logger.Info(sm.Name, fields...)
}

func (r *Reporter) Close() error {
	return nil
}
