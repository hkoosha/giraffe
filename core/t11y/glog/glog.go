package glog

import (
	"errors"
	"sync/atomic"
)

var (
	once          = atomic.Bool{}
	errAlreadySet = errors.New("global logger already set")
)

func Global() Lg {
	return global
}

func SetGlobal(lg Lg) error {
	if !once.CompareAndSwap(false, true) {
		return errAlreadySet
	}

	global = lg

	return nil
}

type Lg interface {
	Named(string) Lg

	Log(level Level, msg string, fields ...any)
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)

	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
	IsError() bool
}
