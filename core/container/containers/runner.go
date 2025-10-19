package containers

import (
	"sync"

	"github.com/hkoosha/giraffe/core/t11y"
	"github.com/hkoosha/giraffe/core/t11y/glog"
	"github.com/hkoosha/giraffe/core/t11y/gtx"

	"github.com/hkoosha/giraffe/core/container/containers/internal"
)

const ListenO11y = "127.0.0.1:8081"

type Runner interface {
	internal.Sealed

	Cycle(gtx.Context, ...Container) error
	MustCycle(gtx.Context, ...Container)

	// =====================================

	Open(gtx.Context) glog.Lg

	Register(...Container)

	Finalize(gtx.Context)

	Wait(gtx.Context) error

	MustWait(gtx.Context)

	Stop(gtx.Context) error

	MustStop(gtx.Context)

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
	t11y.NonNil(ctx, cfg)

	return &runner{
		Sealer:     internal.Sealer{},
		lg:         nil,
		state:      stateWaitingOpen,
		mu:         &sync.Mutex{},
		containers: make([]Container, 0),
		cfg:        cfg,
	}
}

func Configured(
	ctx gtx.Context,
	appRef string,
) Runner {
	cfg := Configure(appRef)
	return GiraffeRunner(ctx, cfg)
}
