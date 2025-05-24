package gtx

import (
	"context"
	rand1 "math/rand"
	rand2 "math/rand/v2"
	"time"
)

type Rand interface {
	StdV1() *rand1.Rand
	StdV2() *rand2.Rand
}

type Clock interface {
	Now() time.Time
}

type Context interface {
	context.Context

	Rand() Rand
	Clock() Clock

	With(k, v any) Context
	WithTimeout(time.Duration) (Context, context.CancelFunc)

	Group() (Context, Group)

	Event(any)
	Debug() []string
}

type Group interface {
	Wait() error
	Go(f func() error)
	SetLimit(int)
	TryGo(f func() error) bool
}

func Of(ctx context.Context) Context {
	if gtx, ok := extract(ctx); ok {
		return gtx
	}

	return set(ctx)
}
