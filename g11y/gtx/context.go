package gtx

import (
	"context"
	rand1 "math/rand"
	rand2 "math/rand/v2"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Baggage interface {
	Span(
		name string,
		opts ...trace.SpanStartOption,
	) (Context, trace.Span)

	StdRandV1() *rand1.Rand

	StdRandV2() *rand2.Rand
}

type Context interface {
	context.Context
	Baggage

	WithTimeout(time.Duration) (Context, context.CancelFunc)
	WithDeadline(time.Time) (Context, context.CancelFunc)
}

func Of(
	parent context.Context,
) Context {
	if gtx, ok := parent.(*giraffeCtx); ok {
		return gtx
	}

	return &giraffeCtx{
		parent:  parent,
		baggage: mkBaggage(parent),
	}
}
