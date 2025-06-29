package containers

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/g11y/containers/internal"
	"github.com/hkoosha/giraffe/g11y/glog"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

const (
	stateWaitingOpen     = "waiting_open"
	stateWaitingFinalize = "waiting_finalize"
	stateWaitingActive   = "waiting_active"
	stateActive          = "active"
	stateClosed          = "closed"
)

// TODO broken implementation.
type runner struct {
	internal.Sealer

	lg glog.Lg

	state      string
	cfg        Config
	mu         *sync.Mutex
	containers []Container
}

func (r *runner) goTo(
	from string,
	to string,
) {
	if from == to {
		panic(EF("cannot go from same state to itself, current=%v, from=to=%v", r.state, from))
	}
	if r.state != from {
		panic(EF("invalid state transition, current=%v, from=%v, to=%v", r.state, from, to))
	}

	r.state = to
}

// Open
// TODO: log, otel.
func (r *runner) Open(
	context.Context,
) glog.Lg {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.goTo(stateWaitingOpen, stateWaitingFinalize)

	return r.lg
}

func (r *runner) Register(
	c ...Container,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != stateWaitingFinalize {
		panic(EF("invalid state, current=%v, expecting=%v", r.state, stateWaitingFinalize))
	}

	r.containers = append(r.containers, c...)
}

// Finalize
// TODO: o11y.Finalize(ctx).
func (r *runner) Finalize(
	ctx context.Context,
	c ...Container,
) {
	r.Register(c...)

	r.mu.Lock()
	defer r.mu.Unlock()

	r.goTo(stateWaitingFinalize, stateWaitingActive)

	if len(r.containers) == 0 {
		panic(EF("no containers registered"))
	}

	for _, c := range r.containers {
		c.Open(ctx, r.lg)
	}

	r.goTo(stateWaitingActive, stateActive)
}

func (r *runner) Wait(
	ctx context.Context,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.goTo(stateWaitingActive, stateActive)

	wg, ctx := errgroup.WithContext(ctx)
	for _, c := range r.containers {
		wg.Go(func() error {
			return c.Start(ctx)
		})
	}

	err := wg.Wait()
	g11y.DieIf(err)
}

func (r *runner) Close(
	ctx context.Context,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.goTo(stateActive, stateClosed)

	ch := make(chan error)

	for _, c := range r.containers {
		launch(ctx, ch, c.Stop)
	}

	select {
	case err := <-ch:
		var err2 any = err

		g11y.DieIf(err2)
	case <-time.After(200 * time.Millisecond):
	}

	// TODO
	// o11y.Shutdown().
}

func launch(
	ctx context.Context,
	ch chan error,
	fn func(context.Context),
) {
	go func() {
		defer func() {
			g11y.DieIf(recover())
		}()

		fn(ctx)

		ch <- nil
	}()
}
