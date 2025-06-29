package containers

import (
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/hkoosha/giraffe/g11y/containers/internal"
	"github.com/hkoosha/giraffe/g11y/glog"
	"github.com/hkoosha/giraffe/g11y/gtx"
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
type runner[D any] struct {
	internal.Sealer

	lg glog.Lg

	state      string
	mu         *sync.Mutex
	containers []Container[D]
}

func (r *runner[D]) goFromTo(
	from string,
	to string,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != from || from == to {
		panic(EF("invalid transition: [%s=>%s]", r.state, to))
	}

	r.state = to
}

func (r *runner[D]) goTo(
	to string,
) {
	// It's still possible r.state could be modified,
	// But I'm not gonna prevent make sure it won't be.
	r.goFromTo(r.state, to)
}

func (r *runner[D]) mustBeIn(
	state string,
) {
	if r.state != state {
		panic(EF("invalid state, current=%v, expecting=%v", r.state, state))
	}
}

// Open
// TODO: log, otel.
func (r *runner[D]) Open(
	_ gtx.Context,
	lg glog.Lg,
) {
	r.goFromTo(stateWaitingOpen, stateWaitingFinalize)
	r.lg = lg
}

func (r *runner[D]) Register(
	c ...Container[D],
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.mustBeIn(stateWaitingFinalize)

	r.containers = append(r.containers, c...)
}

// Finalize
// TODO: o11y.Finalize(ctx).
func (r *runner[D]) Finalize(
	ctx gtx.Context,
	dependencies D,
) {
	r.goFromTo(stateWaitingFinalize, stateWaitingActive)

	if len(r.containers) == 0 {
		panic(EF("no containers registered"))
	}

	for _, c := range r.containers {
		c.Open(ctx, r.lg, dependencies)
	}

	r.goFromTo(stateWaitingActive, stateActive)
}

func (r *runner[D]) Wait(
	ctx gtx.Context,
) error {
	r.goFromTo(stateActive, stateTryingRunning)

	wg, _ctx := errgroup.WithContext(ctx)
	ctx = gtx.Of(_ctx)
	for _, c := range r.containers {
		wg.Go(func() error {
			return c.Run(ctx)
		})
	}

	r.goFromTo(stateTryingRunning, stateRunning)

	return wg.Wait()
}

func (r *runner[D]) MustWait(
	ctx gtx.Context,
) {
	if err := r.Wait(ctx); err != nil {
		panic(E(err))
	}
}

// Stop
// TODO timeout.
func (r *runner[D]) Stop(
	ctx gtx.Context,
	timeout time.Duration,
) error {
	_ = timeout

	r.goFromTo(stateRunning, stateStopping)

	var wg errgroup.Group
	for _, c := range r.containers {
		wg.Go(func() error {
			return c.Stop(ctx)
		})
	}

	if err := wg.Wait(); err != nil {
		r.goFromTo(stateStopping, stateErr)
		return err
	}

	return nil
}

// Close
// TODO implement
// TODO: o11y.Shutdown().
// TODO timeout.
func (r *runner[D]) Close(
	ctx gtx.Context,
	timeout time.Duration,
) {
	_ = timeout

	r.goTo(stateClosing)

	for _, c := range r.containers {
		c.Close(ctx)
	}

	r.goTo(stateClosed)
}
