package containers

import (
	"github.com/hkoosha/giraffe/core/t11y/glog"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
)

type SimpleContainer interface {
	Open(gtx.Context, glog.Lg)

	Start(gtx.Context) error
}

type Container interface {
	Open(gtx.Context, glog.Lg)

	Start(gtx.Context) error

	Stop(gtx.Context) error

	Close(gtx.Context) error
}

func Of(
	fn func(gtx.Context, glog.Lg) error,
) Container {
	return &lgContainer{
		lg: nil,
		fn: fn,
	}
}

func OfMust(
	fn func(gtx.Context, glog.Lg),
) Container {
	return Of(func(ctx gtx.Context, lg glog.Lg) error {
		fn(ctx, lg)
		return nil
	})
}

func FromSimple(
	sc SimpleContainer,
) Container {
	return &simpleAdapter{
		simple: sc,
	}
}
