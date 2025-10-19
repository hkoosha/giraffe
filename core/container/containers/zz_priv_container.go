package containers

import (
	"github.com/hkoosha/giraffe/core/t11y/glog"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
)

var _ Container = (*lgContainer)(nil)

type lgContainer struct {
	lg glog.Lg
	fn func(gtx.Context, glog.Lg) error
}

func (c *lgContainer) Open(_ gtx.Context, lg glog.Lg) {
	c.lg = lg
	c.lg.Debug("open")
}

func (c *lgContainer) Start(ctx gtx.Context) error {
	return c.fn(ctx, c.lg)
}

func (c *lgContainer) Stop(gtx.Context) error {
	c.lg.Debug("stop")
	return nil
}

func (c *lgContainer) Close(gtx.Context) error {
	c.lg.Debug("close")
	return nil
}

// =============================================================================

var _ Container = (*simpleAdapter)(nil)

type simpleAdapter struct {
	simple SimpleContainer
}

func (c *simpleAdapter) Open(ctx gtx.Context, lg glog.Lg) {
	c.simple.Open(ctx, lg)
}

func (c *simpleAdapter) Start(ctx gtx.Context) error {
	return c.simple.Start(ctx)
}

func (c *simpleAdapter) Stop(gtx.Context) error {
	return nil
}

func (c *simpleAdapter) Close(gtx.Context) error {
	return nil
}
