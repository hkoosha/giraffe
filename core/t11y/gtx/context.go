package gtx

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hkoosha/giraffe/core/container/setup"
	"golang.org/x/sync/errgroup"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

var gtxKey = setup.Key("giraffe_core_gtx")

var _ context.Context = (*impl)(nil)

//nolint:containedctx
type impl struct {
	ctx    context.Context
	events *events
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

func (c impl) With(k, v any) Context {
	return &impl{
		ctx:    context.WithValue(c.ctx, k, v),
		events: c.events,
	}
}

func (c impl) WithTimeout(d time.Duration) (Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.ctx, d)

	return &impl{
		ctx:    ctx,
		events: c.events,
	}, cancel
}

func (c impl) Event(v any) {
	c.events.add(v)
}

func (c impl) Debug() []string {
	store := c.events.get()
	s := make([]string, len(store))
	for i := range store {
		s[i] = fmt.Sprintf("%v", store[i])
	}
	return s
}

func (c impl) Group() (Context, Group) {
	group, ctx := errgroup.WithContext(c)

	return &impl{
		ctx:    ctx,
		events: c.events,
	}, group
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

func extract(
	ctx context.Context,
) (*impl, bool) {
	if gtx, ok := ctx.Value(gtxKey).(*impl); ok {
		return gtx, true
	}

	return nil, false
}

func set(
	ctx context.Context,
) *impl {
	if _, ok := extract(ctx); ok {
		panic(EF("gtx already set"))
	}

	gtx := &impl{
		ctx: nil,
		events: &events{
			mu:    &sync.Mutex{},
			store: []any{},
		},
	}

	//nolint:fatcontext
	gtx.ctx = context.WithValue(ctx, gtxKey, gtx)

	return gtx
}
