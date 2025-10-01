package gtx

import (
	"context"

	"github.com/hkoosha/giraffe/g11y/gtx/internal"
)

type Context interface {
	context.Context
}

func Of(ctx context.Context) Context {
	if gtx, ok := internal.Extract(ctx); ok {
		return gtx
	}

	return internal.Set(ctx)
}
