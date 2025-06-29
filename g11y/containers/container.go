package containers

import (
	"context"

	"github.com/hkoosha/giraffe/g11y/glog"
)

type Container interface {
	Open(context.Context, glog.Lg)

	Start(context.Context) error

	Stop(context.Context) error

	Close(context.Context) error
}
