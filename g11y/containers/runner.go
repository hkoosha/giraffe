package containers

import (
	"sync"

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/g11y/containers/internal"
	"github.com/hkoosha/giraffe/g11y/glog"
	"github.com/hkoosha/giraffe/g11y/gtx"
)

const ListenO11y = "127.0.0.1:8081"

type Runner interface {
	internal.Sealed

	Open(ctx gtx.Context) glog.Lg

	Register(...Container)

	Finalize(gtx.Context, ...Container)

	Wait(gtx.Context) error

	MustWait(gtx.Context)

	Stop(ctx gtx.Context) error

	Close(ctx gtx.Context)
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
	ctx gtx.Context,
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
