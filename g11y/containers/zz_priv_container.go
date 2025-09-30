package containers

import (
	"context"

	"github.com/hkoosha/giraffe/g11y/glog"
)

var _ Container = (*lgContainer)(nil)

type lgContainer struct {
	lg glog.Lg
	fn func(context.Context, glog.Lg) error
}

func (c *lgContainer) Open(_ context.Context, lg glog.Lg) {
	c.lg = lg
	c.lg.Debug("open")
}

func (c *lgContainer) Start(ctx context.Context) error {
	return c.fn(ctx, c.lg)
}

func (c *lgContainer) Stop(context.Context) error {
	c.lg.Debug("stop")
	return nil
}

func (c *lgContainer) Close(context.Context) error {
	c.lg.Debug("close")
	return nil
}

// =============================================================================

var _ Container = (*simpleAdapter)(nil)

type simpleAdapter struct {
	simple SimpleContainer
}

func (c *simpleAdapter) Open(ctx context.Context, lg glog.Lg) {
	c.simple.Open(ctx, lg)
}

func (c *simpleAdapter) Start(ctx context.Context) error {
	return c.simple.Start(ctx)
}

func (c *simpleAdapter) Stop(context.Context) error {
	return nil
}

func (c *simpleAdapter) Close(context.Context) error {
	return nil
}
