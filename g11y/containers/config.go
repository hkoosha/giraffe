package containers

import (
	"context"

	"github.com/hkoosha/giraffe/g11y/containers/internal"
	"github.com/hkoosha/giraffe/t11y/glog"
)

type DebugCfg interface {
	internal.Sealed

	IsDebug() bool
}

type LgCfg interface {
	internal.Sealed

	IsLogHumanReadable() bool
	GetLgLevel() glog.Level
}

type OtelCfg interface {
	internal.Sealed

	IsOtel() bool

	GetListenO11y() string
}

type Config interface {
	internal.Sealed

	GetAppRef() string

	Runner(context.Context) Runner
	Wait(context.Context, ...Container) error
	WaitOrDie(context.Context, ...Container)

	DebugCfg
	LgCfg
	OtelCfg
}

// ====================================.

type DebugCfgWrite interface {
	internal.Sealed

	WithDebug() ConfigWrite
	WithoutDebug() ConfigWrite
	SetDebug(bool) ConfigWrite
}

type LgCfgWrite interface {
	internal.Sealed

	WithLgLevel(glog.Level) ConfigWrite

	WithLogHumanReadable() ConfigWrite
	WithoutLogHumanReadable() ConfigWrite
	SetLogHumanReadable(bool) ConfigWrite
}

type OtelCfgWrite interface {
	internal.Sealed

	WithOtel() ConfigWrite
	WithoutOtel() ConfigWrite
	SetOtel(bool) ConfigWrite

	WithListenO11y(string) ConfigWrite
}

type ConfigWrite interface {
	internal.Sealed

	Config

	WithAppRef(string) ConfigWrite

	DebugCfgWrite
	LgCfgWrite
	OtelCfgWrite
}
