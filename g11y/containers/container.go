package containers

import (
	"github.com/hkoosha/giraffe/g11y/glog"
	"github.com/hkoosha/giraffe/g11y/gtx"
)

const (
	StateInit    State = "init"
	StateOpened  State = "opened"
	StateRunning State = "running"
	StateStopped State = "stopped"
	StateClosed  State = "closed"
	StateInvalid State = "invalid"
)

type State string

type Container[D any] interface {
	Open(gtx.Context, glog.Lg, D)

	Run(gtx.Context) error

	Stop(gtx.Context) error

	Close(gtx.Context)

	GetState() State
}
