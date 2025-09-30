package containers

import (
	"context"
	"slices"
	"sync"

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
	stateTryingRunning   = "trying_running"
	stateRunning         = "running"
	stateStopping        = "stopping"
	stateClosing         = "closing"
	stateClosed          = "closed"
	stateErr             = "err"
)

// TODO broken implementation.
type runner struct {
	internal.Sealer

	lg glog.Lg

	state      string
	cfg        Config //nolint:unused
	mu         *sync.Mutex
	containers []Container
}

func (r *runner) goToFrom(
	to string,
	from ...string,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if slices.Contains(from, to) {
		panic(EF("cannot go from same state to itself, current=%v, "+
			"from%v, to=%v",
			r.state,
			from,
			to))
	}
	if slices.Contains(from, r.state) {
		panic(EF("invalid state transition, current=%v, from=%v, to=%v",
			r.state, from, to))
	}

	r.state = to
}

func (r *runner) goTo(
	from string,
	to string,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if from == to {
		panic(EF("cannot go from same state to itself, current=%v, from=to=%v", r.state, from))
	}
	if r.state != from {
		panic(EF("invalid state transition, current=%v, from=%v, to=%v", r.state, from, to))
	}

	r.state = to
}

// Return dummy error, so makes it clear there's an unlock fn that should be
// called via defer. By not returning error, the call looks
// like `defer r.stayIn(foo)()` while the second parentheses are easy to miss.
func (r *runner) stayIn(
	state string,
) (func(), error) {
	r.mu.Lock()
	r.mustBeIn(state)
	return r.mu.Unlock, nil
}

func (r *runner) mustBeIn(
	state string,
) {
	if r.state != state {
		panic(EF("invalid state, current=%v, expecting=%v", r.state, state))
	}
}

// Open
// TODO: log, otel.
func (r *runner) Open(
	context.Context,
) glog.Lg {
	r.goTo(stateWaitingOpen, stateWaitingFinalize)

	panic("todo: set r.lg")
}

func (r *runner) Register(
	c ...Container,
) {
	defer M(r.stayIn(stateWaitingFinalize))

	r.containers = append(r.containers, c...)
}

// Finalize
// TODO: o11y.Finalize(ctx).
func (r *runner) Finalize(
	ctx context.Context,
	c ...Container,
) {
	r.Register(c...)

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
) error {
	r.goToFrom(stateActive, stateTryingRunning)

	wg, ctx := errgroup.WithContext(ctx)
	for _, c := range r.containers {
		wg.Go(func() error {
			return c.Start(ctx)
		})
	}

	r.goToFrom(stateTryingRunning, stateRunning)

	return wg.Wait()
}

func (r *runner) MustWait(
	ctx context.Context,
) {
	M(0, r.Wait(ctx))
}

// Stop
// TODO timeout.
func (r *runner) Stop(
	ctx context.Context,
) error {
	r.goTo(stateRunning, stateStopping)

	var wg errgroup.Group
	for _, c := range r.containers {
		wg.Go(func() error {
			return c.Stop(ctx)
		})
	}

	if err := wg.Wait(); err != nil {
		r.goTo(stateStopping, stateErr)
		return err
	}

	return nil
}

// Close
// TODO implement
// TODO: o11y.Shutdown().
func (r *runner) Close(
	ctx context.Context,
) {
	r.goTo(stateActive, stateClosing)

	var err error
	for _, c := range r.containers {
		err = g11y.Join(err, c.Close(ctx))
	}

	g11y.DieIf(err)

	r.goTo(stateClosing, stateClosed)
}
