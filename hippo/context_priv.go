package hippo

import (
	"context"
	randc "crypto/rand"
	"math/big"
	rand1 "math/rand"
	rand2 "math/rand/v2"
	"time"

	. "github.com/hkoosha/giraffe/internal/dot0"
)

var int63 = big.NewInt(int64(1) << 62)

func seed() int64 {
	return M(randc.Int(randc.Reader, int63)).Int64()
}

type rnd struct {
	seed func() int64
}

func (r *rnd) shallow() *rnd {
	return &rnd{
		seed: r.seed,
	}
}

func (r *rnd) StdV1() *rand1.Rand {
	//nolint:gosec
	return rand1.New(rand1.NewSource(seed()))
}

func (r *rnd) StdV2() *rand2.Rand {
	//nolint:gosec
	return rand2.New(rand2.NewPCG(uint64(r.seed()), uint64(r.seed())))
}

// =============================================================================.

func hContextOf(ctx context.Context) *hContext {
	return &hContext{
		ctx: ctx,
		rnd: &rnd{
			seed: seed,
		},
	}
}

//nolint:containedctx
type hContext struct {
	ctx context.Context
	rnd *rnd
}

func (h *hContext) Rand() Rand {
	return h.rnd
}

func (h *hContext) Deadline() (time.Time, bool) {
	return h.ctx.Deadline()
}

func (h *hContext) Done() <-chan struct{} {
	return h.ctx.Done()
}

func (h *hContext) Err() error {
	return h.ctx.Err()
}

func (h *hContext) Value(key any) any {
	return h.ctx.Value(key)
}

func (h *hContext) WithCtx(ctx context.Context) Context {
	return h.withCtx(ctx)
}

func (h *hContext) withCtx(ctx context.Context) *hContext {
	return &hContext{
		ctx: ctx,
		rnd: h.rnd.shallow(),
	}
}
