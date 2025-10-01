package hippo

import (
	"context"
	rand1 "math/rand"
	rand2 "math/rand/v2"

	"github.com/hkoosha/giraffe/g11y/gtx"
)

type Context interface {
	gtx.Context

	WithCtx(context.Context) Context

	Rand() Rand
}

type Rand interface {
	StdV1() *rand1.Rand
	StdV2() *rand2.Rand
}
