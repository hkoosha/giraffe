package glog

import (
	"github.com/hkoosha/giraffe/g11y/setup"
)

func Global() Lg {
	return global
}

func SetGlobal(lg Lg) {
	setup.Finish("giraffe", "log", "global")
	global = lg
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
