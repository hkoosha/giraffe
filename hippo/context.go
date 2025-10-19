package hippo

import (
	"context"

	"github.com/hkoosha/giraffe/core/t11y/gtx"
)

type Context = gtx.Context

func ContextOf(ctx context.Context) gtx.Context {
	return gtx.Of(ctx)
}
