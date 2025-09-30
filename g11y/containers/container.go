package containers

import (
	"context"

	"github.com/hkoosha/giraffe/g11y/glog"
)

type SimpleContainer interface {
	Open(context.Context, glog.Lg)

	Start(context.Context) error
}

type Container interface {
	Open(context.Context, glog.Lg)

	Start(context.Context) error

	Stop(context.Context) error

	Close(context.Context) error
}

func Of(
	fn func(context.Context, glog.Lg) error,
) Container {
	return &lgContainer{
		fn: fn,
	}
}

func OfMust(
	fn func(context.Context, glog.Lg),
) Container {
	return Of(func(ctx context.Context, lg glog.Lg) error {
		fn(ctx, lg)
		return nil
	})
}

func FromSimple(sc SimpleContainer) Container {
	return &lgContainer{}
}
