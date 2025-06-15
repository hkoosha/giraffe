package gzapadapter

import (
	"strconv"

	"go.uber.org/zap"

	"github.com/hkoosha/giraffe/glog"
)

func Of(lg *zap.Logger) glog.Lg {
	return Hooked(
		lg,
		nil,
		nil,
		nil,
		nil,
	)
}

func Hooked(
	lg *zap.Logger,
	debugHook Hook,
	infoHook Hook,
	warnHook Hook,
	errorHook Hook,
) glog.Lg {
	return &adapter{
		lg:        lg,
		debugHook: debugHook,
		infoHook:  infoHook,
		warnHook:  warnHook,
		errorHook: errorHook,
	}
}

// ============================================================================.

type adapter struct {
	lg        *zap.Logger
	debugHook Hook
	infoHook  Hook
	warnHook  Hook
	errorHook Hook
}

func (z adapter) Named(s string) glog.Lg {
	return Hooked(
		z.lg.Named(s),
		z.debugHook,
		z.infoHook,
		z.warnHook,
		z.errorHook,
	)
}

func (z adapter) Debug(msg string, fields ...any) {
	z.lg.Debug(msg, toZap("", fields)...)
}

func (z adapter) Info(msg string, fields ...any) {
	z.lg.Info(msg, toZap("", fields)...)
}

func (z adapter) Warn(msg string, fields ...any) {
	z.lg.Warn(msg, toZap("", fields)...)
}

func (z adapter) Error(msg string, fields ...any) {
	z.lg.Error(msg, toZap("", fields)...)
}

func (z adapter) Err(msg string, err error, fields ...any) {
	z.lg.Error(msg, toZap("", append(fields, err))...)
}

func (z adapter) Of(key string, value ...any) any {
	return toZap(key, value)
}

func toZap(
	key string,
	fields []any,
) []zap.Field {
	hasKey := key != ""

	list := make([]zap.Field, len(fields))
	for i, f := range fields {
		switch v := f.(type) {
		case zap.Field:
			list[i] = v

		case glog.KV:
			list[i] = zap.Any(v.Key(), v.Val())

		case error:
			switch {
			case hasKey && i > 0:
				list[i] = zap.NamedError(key+"_"+strconv.Itoa(i), v)

			case hasKey:
				list[i] = zap.NamedError(key, v)

			// case hasIndex:
			//   Handled by default case.
			//   TODO what does zap do on multiple unnamed err?

			default:
				list[i] = zap.Error(v)
			}

		default:
			var vKey string
			switch {
			case hasKey && i > 0:
				vKey = key + "_" + strconv.Itoa(i)

			case i > 0:
				vKey = "f_" + strconv.Itoa(i)

			case hasKey:
				vKey = key

			default:
				vKey = "f0"
			}

			list[i] = zap.Any(vKey, v)
		}
	}

	return list
}
