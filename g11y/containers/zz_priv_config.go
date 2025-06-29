package containers

import (
	"context"

	"github.com/hkoosha/giraffe/g11y/containers/internal"
	"github.com/hkoosha/giraffe/g11y/glog"
)

type config struct {
	internal.Sealer
	appRef        string
	listenO11y    string
	otelEndpoint  string
	level         glog.Level
	debug         bool
	humanReadable bool
	otel          bool
	otelInsecure  bool
}

func (r *config) shallow() *config {
	return &*r
}

// ============================================================================.

func (r *config) Runner(
	ctx context.Context,
) Runner {
	return GiraffeRunner(ctx, r)
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

func (r *config) WithOtelEndpoint(s string) ConfigWrite {
	cp := r.shallow()
	cp.otelEndpoint = s
	return cp
}

func (r *config) WithOtelInsecure() ConfigWrite {
	return r.SetOtelInsecure(true)
}

func (r *config) WithoutOtelInsecure() ConfigWrite {
	return r.SetOtelInsecure(false)
}

func (r *config) SetOtelInsecure(b bool) ConfigWrite {
	cp := r.shallow()
	cp.otelInsecure = b
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

func (r *config) GetOtelEndpoint() string {
	return r.otelEndpoint
}

func (r *config) IsOtelInsecure() bool {
	return r.otelInsecure
}

func (r *config) GetListenO11y() string {
	return r.listenO11y
}
