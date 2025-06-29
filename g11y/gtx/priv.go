package gtx

import (
	"context"
	randc "crypto/rand"
	"math/big"
	rand1 "math/rand"
	rand2 "math/rand/v2"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/hkoosha/giraffe/g11y/gotel"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

type baggage struct {
	seed   func() int64
	tracer trace.Tracer
}

func mkBaggage(
	ctx context.Context,
) *baggage {
	switch v := ctx.Value(baggageKey).(type) {
	case nil:
		b := newBaggage()
		return b

	case *baggage:
		return v

	default:
		panic(EF(
			"unreachable unknown baggage type: %T",
			ctx.Value(baggageKey),
		))
	}
}

func newBaggage() *baggage {
	return &baggage{
		tracer: gotel.Tracer(),
		seed: func() int64 {
			return M(randc.Int(randc.Reader, int63)).Int64()
		},
	}
}

type baggageKeyT struct{}

var baggageKey baggageKeyT

// ============================================================================.

var _ Context = (*giraffeCtx)(nil)

//nolint:containedctx
type giraffeCtx struct {
	parent  context.Context
	baggage *baggage
}

//nolint:spancheck
func (g giraffeCtx) Span(
	name string,
	opts ...trace.SpanStartOption,
) (Context, trace.Span) {
	ctx, span := g.baggage.tracer.Start(g, name, opts...)
	return Of(ctx), span
}

func (g giraffeCtx) Deadline() (time.Time, bool) {
	return g.parent.Deadline()
}

func (g giraffeCtx) Done() <-chan struct{} {
	return g.parent.Done()
}

func (g giraffeCtx) Err() error {
	return g.parent.Err()
}

func (g giraffeCtx) Value(v any) any {
	return g.parent.Value(v)
}

// ====================================.

var int63 = big.NewInt(int64(1) << 62)

func seed() int64 {
	return M(randc.Int(randc.Reader, int63)).Int64()
}

func (g giraffeCtx) StdRandV1() *rand1.Rand {
	//nolint:gosec
	return rand1.New(rand1.NewSource(seed()))
}

func (g giraffeCtx) StdRandV2() *rand2.Rand {
	//nolint:gosec
	return rand2.New(rand2.NewPCG(uint64(g.baggage.seed()), uint64(g.baggage.seed())))
}

// ====================================.

func (g giraffeCtx) WithTimeout(
	timeout time.Duration,
) (Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(g.parent, timeout)
	return Of(ctx), cancel
}

func (g giraffeCtx) WithDeadline(
	deadline time.Time,
) (Context, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(g.parent, deadline)
	return Of(ctx), cancel
}
