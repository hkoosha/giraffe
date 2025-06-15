package glog

import (
	"time"

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

	Err(msg string, err error, fields ...any)

	Of(key string, value ...any) any
}

type KVSafe interface {
	bool |
		string |
		[]byte |
		float32 |
		float64 |
		int |
		int8 |
		int16 |
		int32 |
		int64 |
		uint |
		uint8 |
		uint16 |
		uint32 |
		uint64 |
		time.Time |
		time.Duration |
		*bool |
		*string |
		*[]byte |
		*float32 |
		*float64 |
		*int |
		*int8 |
		*int16 |
		*int32 |
		*int64 |
		*uint |
		*uint8 |
		*uint16 |
		*uint32 |
		*uint64 |
		*time.Time |
		*time.Duration
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

func Of[V KVSafe](
	key string,
	v V,
) KV {
	return KV{key: key, val: v}
}
