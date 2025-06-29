package containers

import (
	"context"
	"sync"
	"time"

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/g11y/containers/internal"
	"github.com/hkoosha/giraffe/g11y/glog"
)

type Runner interface {
	internal.Sealed

	Open(ctx context.Context) glog.Lg

	Register(...Container)

	Finalize(context.Context, ...Container)

	Wait(context.Context) error

	Stop(ctx context.Context, timeout time.Duration) error

	Close(ctx context.Context, timeout time.Duration)
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

func GiraffeRunner(
	ctx context.Context,
	cfg ConfigWrite,
) Runner {
	g11y.NonNil(ctx, cfg)

	return &runner{
		Sealer: internal.Sealer{},

		state:      stateWaitingOpen,
		mu:         &sync.Mutex{},
		containers: make([]Container, 0),
		cfg:        cfg,
	}
}
