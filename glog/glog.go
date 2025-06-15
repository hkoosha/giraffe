package glog

import (
	"github.com/hkoosha/giraffe/g11y/setup"
)

type Lg interface {
	Named(string) Lg

	Debug(msg string, fields ...any)

	Info(msg string, fields ...any)

	Warn(msg string, fields ...any)

	Error(msg string, fields ...any)

	Err(msg string, err error, fields ...any)

	Of(key string, value ...any) any
}

func Global() Lg {
	return global
}

func SetGlobal(lg Lg) {
	setup.Once("giraffe", "log", "global")
	global = lg
}
