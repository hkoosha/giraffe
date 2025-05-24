package containers

import (
	"slices"
	"sync"

	"github.com/hkoosha/giraffe/contrib/zap/gzapadapter"
	"github.com/hkoosha/giraffe/core/t11y"
	"github.com/hkoosha/giraffe/core/t11y/glog"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	"golang.org/x/sync/errgroup"

	"github.com/hkoosha/giraffe/core/container/containers/internal"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

const (
	stateWaitingOpen     = "waiting_open"
	stateWaitingFinalize = "waiting_finalize"
	stateWaitingActive   = "waiting_active"
	ready                = "ready"
	stateTryingRunning   = "trying_running"
	stateRunning         = "running"
	stateStopping        = "stopping"
	stateStopped         = "stopped"
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

func (r *runner) gotoFrom(
	to string,
	from ...string,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if slices.Contains(from, to) {
		panic(EF("cannot go from a state to itself, current=%v, "+
			"from%v, to=%v",
			r.state,
			from,
			to))
	}
	if !slices.Contains(from, r.state) {
		panic(EF("invalid state transition, current=%v, from=%v, to=%v",
			r.state, from, to))
	}

	r.state = to
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
	gtx.Context,
) glog.Lg {
	r.gotoFrom(stateWaitingFinalize, stateWaitingOpen)

	lgProvider := gzapadapter.MkInit(true, true, nil)
	M(lgProvider.Open())

	r.lg = gzapadapter.Of(lgProvider.Get())
	return r.lg
}

func (r *runner) Register(
	c ...Container,
) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mustBeIn(stateWaitingFinalize)

	r.containers = append(r.containers, c...)
}

// Finalize
// TODO: o11y.Finalize(ctx).
func (r *runner) Finalize(
	ctx gtx.Context,
) {
	r.gotoFrom(stateWaitingActive, stateWaitingFinalize)

	if len(r.containers) == 0 {
		panic(EF("no containers registered"))
	}

	for _, c := range r.containers {
		c.Open(ctx, r.lg)
	}

	r.gotoFrom(ready, stateWaitingActive)
}

func (r *runner) Wait(
	ctx gtx.Context,
) error {
	r.gotoFrom(stateTryingRunning, ready)

	ctx, wg := ctx.Group()
	for _, c := range r.containers {
		wg.Go(func() error {
			return c.Start(ctx)
		})
	}

	r.gotoFrom(stateRunning, stateTryingRunning)

	return wg.Wait()
}

func (r *runner) MustWait(
	ctx gtx.Context,
) {
	t11y.DieIf(r.Wait(ctx))
}

// Stop
// TODO timeout.
func (r *runner) Stop(
	ctx gtx.Context,
) error {
	r.gotoFrom(stateStopping, stateRunning)

	var wg errgroup.Group
	for _, c := range r.containers {
		wg.Go(func() error {
			return c.Stop(ctx)
		})
	}

	if err := wg.Wait(); err != nil {
		r.gotoFrom(stateErr, stateStopping)
		return err
	}

	r.gotoFrom(stateStopped, stateStopping)
	return nil
}

func (r *runner) MustStop(
	ctx gtx.Context,
) {
	t11y.DieIf(r.Stop(ctx))
}

// Close
// TODO implement
// TODO: o11y.Shutdown().
func (r *runner) Close(
	ctx gtx.Context,
) {
	r.gotoFrom(stateClosing, stateStopped)

	var err error
	for _, c := range r.containers {
		err = t11y.Join(err, c.Close(ctx))
	}

	t11y.DieIf(err)

	r.gotoFrom(stateClosed, stateClosing)
}

func (r *runner) Cycle(
	ctx gtx.Context,
	c ...Container,
) error {
	if len(c) == 0 {
		return EF("must provide at least one container")
	}

	r.Open(ctx)
	r.Register(c...)
	r.Finalize(ctx)

	if err := r.Wait(ctx); err != nil {
		return err
	}

	if err := r.Stop(ctx); err != nil {
		return err
	}

	r.Close(ctx)

	return nil
}

func (r *runner) MustCycle(
	ctx gtx.Context,
	c ...Container,
) {
	t11y.DieIf(r.Cycle(ctx, c...))
}
