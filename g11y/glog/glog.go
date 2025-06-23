package glog

import (
	"github.com/hkoosha/giraffe/g11y/setup"
)

func Global() Lg {
	return global
}

func SetGlobal(lg Lg) {
	setup.Once("giraffe", "log", "global")
	global = lg
}

// ============================================================================.

type Lg interface {
	Named(string) Lg

	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)

	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
	IsError() bool
}

type KV struct {
	val any
	key string
}

func (k *KV) Key() string {
	return k.key
}

func (k *KV) Val() any {
	return k.val
}

func Of(
	key string,
	v any,
) KV {
	return KV{key: key, val: v}
}
