package containers

import (
	"sync"
	"time"

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/g11y/containers/internal"
	"github.com/hkoosha/giraffe/g11y/glog"
	"github.com/hkoosha/giraffe/g11y/gtx"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

func MustRun[D any](
	ctx gtx.Context,
	lg glog.Lg,
	cfg Config,
	dependencies D,
	c Container[D],
) {
	OK(Run[D](ctx, lg, cfg, dependencies, c))
}

func Run[D any](
	ctx gtx.Context,
	lg glog.Lg,
	cfg Config,
	dependencies D,
	c Container[D],
) error {
	r := GiraffeRunner[D](ctx, cfg)
	defer r.Close(ctx, 1*time.Second)

	r.Open(ctx, lg)
	r.Register(c)
	r.Finalize(ctx, dependencies)
	return r.Wait(ctx)
}

type Runner[D any] interface {
	internal.Sealed

	Open(gtx.Context, glog.Lg)

	Register(...Container[D])

	Finalize(gtx.Context, D)

	Wait(gtx.Context) error

	MustWait(gtx.Context)

	Stop(ctx gtx.Context, timeout time.Duration) error

	Close(ctx gtx.Context, timeout time.Duration)
}

// ====================================.

func Configure(
	appRef string,
	listenO11y string,
	otelEndpoint string,
) ConfigWrite {
	return &config{
		Sealer:        internal.Sealer{},
		debug:         false,
		level:         glog.Info,
		humanReadable: false,
		appRef:        appRef,
		otel:          false,
		listenO11y:    listenO11y,
		otelEndpoint:  otelEndpoint,
		otelInsecure:  false,
	}
}

func GiraffeRunner[D any](
	ctx gtx.Context,
	cfg Config,
) Runner[D] {
	_ = cfg

	g11y.NonNil(ctx, cfg)

	return &runner[D]{
		Sealer: internal.Sealer{},

		lg:         nil,
		state:      stateWaitingOpen,
		mu:         &sync.Mutex{},
		containers: make([]Container[D], 0),
	}
}
