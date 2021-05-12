package gotrace

import (
	"fmt"
	"sync"
	"time"

	"github.com/lyouthzzz/go-trace/attribute"
)

// ====SpanKind====
type SpanKind string

const (
	SpanKindProvider SpanKind = "PROVIDER"
	SpanKindConsumer SpanKind = "CUNSUMER"
)

// ====SpanStatus====
type SpanStatus string

const (
	SpanStatusSuccess SpanStatus = "SUCCESS"
	SpanStatusFail    SpanStatus = "FAIL"
)

// ====SpanContext====

type SpanContext struct {
	// 调用链唯一id
	globalTicketId string
	monitorId      string
	parentRpcId    string
	rpcId          string
}

func (c SpanContext) GetGlobalTicketId() string {
	return c.globalTicketId
}

func (c SpanContext) GetMonitorId() string {
	return c.monitorId
}

func (c SpanContext) GetParentRpcId() string {
	return c.parentRpcId
}

func (c SpanContext) GetRpcId() string {
	return c.rpcId
}

func (c SpanContext) IsValid() bool {
	return c.globalTicketId != ""
}

// ====SpanOption====
type SpanOption func(*span)

func SpanKindOption(kind SpanKind) SpanOption {
	return func(s *span) {
		s.kind = kind
	}
}

func SpanTracerOption(tracer Tracer) SpanOption {
	return func(s *span) {
		s.tracer = tracer
	}
}

func SpanContextOption(sc SpanContext) SpanOption {
	return func(s *span) {
		s.spanContext = sc
	}
}

var _ Span = (*span)(nil)

type Span interface {
	// 获取 Tracer
	Tracer() Tracer
	// 保存消息
	End()
	// 获取上下文
	GetSpanContext() SpanContext
	// 添加错误
	AddError(error)
	// 设置name
	SetName(string)
	// 设置span状态
	SetStatus(SpanStatus)
	// 设置属性信息
	SetAttributes(attrs ...attribute.KeyValue)
	// child节点
	Child(name string) Span
	// follow节点
	Follow(name string) Span
}

func NewSpan(name string, opts ...SpanOption) Span {
	s := &span{name: name, startTime: time.Now()}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type span struct {
	lock        sync.Mutex
	tracer      Tracer
	spanContext SpanContext
	name        string
	kind        SpanKind
	status      SpanStatus
	errs        []error
	ended       bool
	startTime   time.Time
	duration    time.Duration
	attributes  map[string]interface{}
	children    int
}

func (s *span) Tracer() Tracer {
	return s.tracer
}

func (s *span) GetSpanContext() SpanContext {
	return s.spanContext
}

func (s *span) End() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	s.duration = time.Since(s.startTime)
	s.ended = true
}

func (s *span) AddError(err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.errs == nil {
		s.errs = make([]error, 0)
	}
	s.errs = append(s.errs, err)
}

func (s *span) SetName(name string) {
	s.name = name
}

func (s *span) SetStatus(status SpanStatus) {
	s.status = status
}

func (s *span) SetAttributes(attrs ...attribute.KeyValue) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	if s.attributes == nil {
		s.attributes = make(map[string]interface{})
	}

	for _, attr := range attrs {
		s.attributes[attr.Key] = attr.Value
	}
}

func (s *span) Child(name string) Span {
	cs := &span{name: name}
	csc := SpanContext{}

	sc := s.GetSpanContext()

	if sc.globalTicketId != "" {
		csc.globalTicketId = sc.globalTicketId
	}
	if sc.rpcId != "" {
		csc.parentRpcId = sc.rpcId
		csc.rpcId = s.childId()
	}
	csc.monitorId = "monitorUUID"

	return cs
}

func (s *span) Follow(name string) Span {
	return &span{}
}

func (s *span) childId() string {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.children++
	return fmt.Sprintf("%s.%d", s.GetSpanContext().GetRpcId(), s.children)
}
