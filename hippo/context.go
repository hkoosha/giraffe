package hippo

import (
	"context"
	rand1 "math/rand"
	rand2 "math/rand/v2"
)

type HContext interface {
	context.Context

	WithCtx(context.Context) HContext

	Rand() Rand
}

type Rand interface {
	StdV1() *rand1.Rand
	StdV2() *rand2.Rand
}
