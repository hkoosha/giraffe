package containers

import (
	"context"
	"sync"

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/g11y/containers/internal"
	"github.com/hkoosha/giraffe/g11y/glog"
)

const ListenO11y = "127.0.0.1:8081"

type Runner interface {
	internal.Sealed

	Open(ctx context.Context) glog.Lg

	Register(...Container)

	Finalize(context.Context, ...Container)

	Wait(context.Context) error

	MustWait(context.Context)

	Stop(ctx context.Context) error

	Close(ctx context.Context)
}

// ====================================.

func Configure(
	appRef string,
) ConfigWrite {
	return &config{
		Sealer:        internal.Sealer{},
		debug:         false,
		level:         glog.Info,
		humanReadable: false,
		appRef:        appRef,
		otel:          false,
		listenO11y:    ListenO11y,
	}
}

func GiraffeRunner(
	ctx context.Context,
	cfg ConfigWrite,
) Runner {
	g11y.NonNil(ctx, cfg)

	return &runner{
		Sealer: internal.Sealer{},

		lg: nil,

		state:      stateWaitingOpen,
		mu:         &sync.Mutex{},
		containers: make([]Container, 0),
		cfg:        cfg,
	}
}
