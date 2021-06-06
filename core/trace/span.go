package trace

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuexiangyou/thor/core/stringx"
	"github.com/xuexiangyou/thor/core/timex"
	"github.com/xuexiangyou/thor/core/trace/tracespec"
)

const (
	initSpanId  = "0"
	clientFlag  = "client"
	serverFlag  = "server"
	spanSeqRune = '.'
)

var spanSep = string([]byte{spanSeqRune})

type Span struct {
	ctx           spanContext
	serviceName   string
	operationName string
	startTime     time.Time
	flag          string
	children      int
}

func newServerSpan(carrier Carrier, serviceName, operationName string) tracespec.Trace {
	traceId := stringx.TakeWithPriority(func() string {
		if carrier != nil {
			return carrier.Get(traceIdKey)
		}
		return ""
	}, stringx.RandId)

	spanId := stringx.TakeWithPriority(func() string {
		if carrier != nil {
			return carrier.Get(spanIdKey)
		}
		return ""
	}, func() string {
		return initSpanId
	})

	return &Span{
		ctx: spanContext{
			traceId: traceId,
			spanId: spanId,
		},
		serviceName: serviceName,
		operationName: operationName,
		startTime: timex.Time(),
		flag: serverFlag,
	}
}

func (s *Span) Finish() {
}

func (s *Span) Follow(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	span := &Span{
		ctx: spanContext{
			traceId: s.ctx.traceId,
			spanId: s.followSpanId(),
		},
		serviceName: serviceName,
		operationName: operationName,
		startTime: timex.Time(),
		flag: s.flag,
	}
	return context.WithValue(ctx, tracespec.TracingKey, span), span
}

func (s *Span) Fork(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	span := &Span{
		ctx: spanContext{
			traceId: s.ctx.traceId,
			spanId: s.forkSpanId(),
		},
		serviceName: serviceName,
		operationName: operationName,
		startTime: timex.Time(),
		flag: clientFlag,
	}
	return context.WithValue(ctx, tracespec.TracingKey, span), span
}

func (s *Span) SpanId() string {
	return s.ctx.spanId
}

func (s *Span) TraceId() string {
	return s.ctx.traceId
}

func (s *Span) Visit(fn func(key, val string) bool) {
	s.ctx.Visit(fn)
}

func (s *Span) forkSpanId() string {
	s.children++
	return fmt.Sprintf("%s.%d", s.ctx.spanId, s.children)
}

func (s *Span) followSpanId() string {
	fields := strings.FieldsFunc(s.ctx.spanId, func(r rune) bool {
		return r == spanSeqRune
	})
	if len(fields) == 0 {
		return s.ctx.spanId
	}

	last := fields[len(fields) - 1]
	val, err := strconv.Atoi(last)
	if err != nil {
		return s.ctx.spanId
	}
	last = strconv.Itoa(val + 1)
	return strings.Join(fields, spanSep)
}

func StartServerSpan(ctx context.Context, carrier Carrier, serviceName, operationName string) (context.Context, tracespec.Trace) {
	span := newServerSpan(carrier, serviceName, operationName)
	return context.WithValue(ctx, tracespec.TracingKey, span), span
}
