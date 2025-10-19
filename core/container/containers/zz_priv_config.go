package containers

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/hkoosha/giraffe/core/t11y"
	"github.com/hkoosha/giraffe/core/t11y/glog"
	"github.com/hkoosha/giraffe/core/t11y/gtx"

	"github.com/hkoosha/giraffe/core/container/containers/internal"
)

var (
	errStopTimeout  = errors.New("timed out on stop")
	errCloseTimeout = errors.New("timed out on close")
)

type config struct {
	internal.Sealer
	appRef        string
	listenO11y    string
	level         glog.Level
	debug         bool
	humanReadable bool
	otel          bool
}

func (r *config) shallow() *config {
	return &*r
}

// ============================================================================.

func (r *config) Runner(
	ctx gtx.Context,
) Runner {
	return GiraffeRunner(ctx, r)
}

func (r *config) Wait(
	ctx gtx.Context,
	containers ...Container,
) error {
	err := atomic.Value{}
	return r.doWait(ctx, &err, containers...)
}

func (r *config) doWait(
	ctx gtx.Context,
	err *atomic.Value,
	containers ...Container,
) error {
	const timeout = 4 * time.Second

	defer func() { t11y.Mix(err, recover()) }()

	open := func() Runner {
		rn := GiraffeRunner(ctx, r)
		rn.Open(ctx)
		rn.Register(containers...)
		rn.Finalize(ctx)
		return rn
	}

	fin := func(
		ctx gtx.Context,
		rn Runner,
	) {
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		ctx, cancel := ctx.WithTimeout(timeout)
		defer cancel()

		done := make(chan error, 1)

		go func() {
			defer func() { t11y.DieIf(recover()) }()
			done <- rn.Stop(ctx)
		}()

		select {
		case dErr := <-done:
			t11y.Mix(err, dErr)
		case <-timer.C:
			t11y.Mix(err, errStopTimeout)
		}

		go func() {
			defer func() { t11y.DieIf(recover()) }()
			rn.Close(ctx)
			done <- nil
		}()

		select {
		case dErr := <-done:
			t11y.Mix(err, dErr)
		case <-timer.C:
			t11y.Mix(err, errCloseTimeout)
		}
	}

	rn := open()
	defer fin(ctx, rn)
	return t11y.MixAndGet(err, rn.Wait(ctx))
}

func (r *config) WaitOrDie(
	ctx gtx.Context,
	containers ...Container,
) {
	err := r.Wait(ctx, containers...)
	t11y.DieIf(err)
}

// ============================================================================.

func (r *config) WithDebug() ConfigWrite {
	return r.SetDebug(true)
}

func (r *config) WithoutDebug() ConfigWrite {
	return r.SetDebug(false)
}

func (r *config) SetDebug(b bool) ConfigWrite {
	cp := r.shallow()
	cp.debug = b
	return cp
}

func (r *config) WithLogHumanReadable() ConfigWrite {
	return r.SetLogHumanReadable(true)
}

func (r *config) WithoutLogHumanReadable() ConfigWrite {
	return r.SetLogHumanReadable(false)
}

func (r *config) SetLogHumanReadable(b bool) ConfigWrite {
	cp := r.shallow()
	cp.humanReadable = b
	return cp
}

func (r *config) WithLgLevel(level glog.Level) ConfigWrite {
	cp := r.shallow()
	cp.level = level
	return cp
}

func (r *config) WithAppRef(s string) ConfigWrite {
	cp := r.shallow()
	cp.appRef = s
	return cp
}

func (r *config) WithOtel() ConfigWrite {
	return r.SetOtel(true)
}

func (r *config) WithoutOtel() ConfigWrite {
	return r.SetOtel(false)
}

func (r *config) SetOtel(b bool) ConfigWrite {
	cp := r.shallow()
	cp.otel = b
	return cp
}

func (r *config) WithListenO11y(s string) ConfigWrite {
	cp := r.shallow()
	cp.listenO11y = s
	return cp
}

// ============================================================================.

func (r *config) IsDebug() bool {
	return r.debug
}

func (r *config) IsLogHumanReadable() bool {
	return r.humanReadable
}

func (r *config) GetLgLevel() glog.Level {
	return r.level
}

func (r *config) GetAppRef() string {
	return r.appRef
}

func (r *config) IsOtel() bool {
	return r.otel
}

func (r *config) GetListenO11y() string {
	return r.listenO11y
}
