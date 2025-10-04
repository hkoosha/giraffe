package gtx

import (
	"context"
	"time"

	. "github.com/hkoosha/giraffe/t11y/dot"
)

type gtxKeyT int

var gtxKey gtxKeyT

// =============================================================================

var _ context.Context = (*impl)(nil)

//nolint:containedctx
type impl struct {
	ctx context.Context
}

func (c impl) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c impl) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c impl) Err() error {
	return c.ctx.Err()
}

func (c impl) Value(key any) any {
	return c.ctx.Value(key)
}

func (c impl) WithCtx(ctx context.Context) Context {
	return &impl{
		ctx: ctx,
	}
}

// =====================================

func (c impl) Clock() Clock {
	return clock{}
}

// =====================================

func (c impl) Rand() Rand {
	return rnd{
		seed: seed,
	}
}

// =============================================================================

func extract(ctx context.Context) (*impl, bool) {
	if gtx, ok := ctx.Value(gtxKey).(*impl); ok {
		return gtx, true
	}

	return nil, false
}

func set(ctx context.Context) *impl {
	if _, ok := extract(ctx); ok {
		panic(EF("gtx already set"))
	}

	gtx := &impl{ctx: nil}
	//nolint:fatcontext
	gtx.ctx = context.WithValue(ctx, gtxKey, gtx)
	return gtx
}
