package gotrace

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lyouthzzz/go-trace/attribute"
)

// ====SpanKind====
type SpanKind string

const (
	SpanKindProvider SpanKind = "provider"
	SpanKindConsumer SpanKind = "consumer"
)

func (kind SpanKind) String() string {
	return string(kind)
}

// ====SpanStatus====
type SpanStatus int

const (
	SpanStatusSuccess SpanStatus = 1
	SpanStatusFail    SpanStatus = 0
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
type SpanOption func(*SpanModel)

func SpanKindOption(kind SpanKind) SpanOption {
	return func(s *SpanModel) {
		s.Kind = kind
	}
}

func SpanTracerOption(t Tracer) SpanOption {
	return func(s *SpanModel) {
		if tracer, ok := t.(*tracer); ok {
			s.tracer = tracer
		}
	}
}

func SpanContextOption(sc SpanContext) SpanOption {
	return func(s *SpanModel) {
		s.SpanContext = sc
	}
}

var _ Span = (*SpanModel)(nil)

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
	s := &SpanModel{Name: name, StartTime: time.Now()}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type SpanModel struct {
	lock        sync.Mutex
	tracer      *tracer
	SpanContext SpanContext
	Name        string
	Kind        SpanKind
	Status      SpanStatus
	Errs        []error
	StartTime   time.Time
	Duration    time.Duration
	Attributes  map[string]interface{}
	ended       bool
	children    int
}

func (s *SpanModel) Tracer() Tracer {
	return s.tracer
}

func (s *SpanModel) GetSpanContext() SpanContext {
	return s.SpanContext
}

func (s *SpanModel) End() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	s.Duration = time.Since(s.StartTime)

	if s.tracer.reporter != nil {
		s.tracer.reporter.Send(s)
	}
	s.ended = true
}

func (s *SpanModel) AddError(err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.Errs == nil {
		s.Errs = make([]error, 0)
	}
	s.Errs = append(s.Errs, err)
}

func (s *SpanModel) SetName(name string) {
	s.Name = name
}

func (s *SpanModel) SetStatus(status SpanStatus) {
	s.Status = status
}

func (s *SpanModel) SetAttributes(attrs ...attribute.KeyValue) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	if s.Attributes == nil {
		s.Attributes = make(map[string]interface{})
	}

	for _, attr := range attrs {
		s.Attributes[attr.Key] = attr.Value
	}
}

func (s *SpanModel) Child(name string) Span {
	cs := &SpanModel{Name: name}
	csc := SpanContext{}

	sc := s.GetSpanContext()

	if sc.globalTicketId != "" {
		csc.globalTicketId = sc.globalTicketId
	}
	if sc.rpcId != "" {
		csc.parentRpcId = sc.rpcId
		csc.rpcId = s.childId()
	}

	csc.monitorId = uuid.New().String()

	return cs
}

func (s *SpanModel) Follow(name string) Span {
	return &SpanModel{}
}

func (s *SpanModel) childId() string {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.children++
	return fmt.Sprintf("%s.%d", s.GetSpanContext().GetRpcId(), s.children)
}
