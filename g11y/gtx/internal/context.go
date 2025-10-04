package internal

import (
	"context"
	"time"

	. "github.com/hkoosha/giraffe/internal/dot0"
)

type gtxKeyT int

var gtxKey gtxKeyT

// =============================================================================

var _ context.Context = (*GContext)(nil)

//nolint:containedctx
type GContext struct {
	ctx context.Context
}

func (c GContext) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c GContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c GContext) Err() error {
	return c.ctx.Err()
}

func (c GContext) Value(key any) any {
	return c.ctx.Value(key)
}

// =============================================================================

func Extract(ctx context.Context) (*GContext, bool) {
	if gtx, ok := ctx.Value(gtxKey).(*GContext); ok {
		return gtx, true
	}

	return nil, false
}

func Set(ctx context.Context) *GContext {
	if _, ok := Extract(ctx); ok {
		panic(EF("gtx already set"))
	}

	gtx := &GContext{ctx: nil}
	//nolint:fatcontext
	gtx.ctx = context.WithValue(ctx, gtxKey, gtx)
	return gtx
}
