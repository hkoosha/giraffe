package gzapadapter

import (
	"strconv"

	"github.com/hkoosha/giraffe/core/t11y/glog"
	"go.uber.org/zap"
)

func Of(lg *zap.Logger) glog.Lg {
	var level glog.Level
	switch {
	case lg.Check(zap.DebugLevel, "") != nil:
		level = glog.Debug

	case lg.Check(zap.InfoLevel, "") != nil:
		level = glog.Info

	case lg.Check(zap.WarnLevel, "") != nil:
		level = glog.Warn

	case lg.Check(zap.ErrorLevel, "") != nil:
		level = glog.Error

	default:
		level = glog.Disabled
	}

	return &adapter{
		lg:  lg,
		max: level,
	}
}

type adapter struct {
	lg  *zap.Logger
	max glog.Level
}

func (z adapter) Named(s string) glog.Lg {
	return &adapter{
		lg:  z.lg.Named(s),
		max: z.max,
	}
}

func (z adapter) Log(level glog.Level, msg string, fields ...any) {
	switch level {
	case glog.Debug:
		z.Debug(msg, fields...)
	case glog.Info:
		z.Info(msg, fields...)
	case glog.Warn:
		z.Warn(msg, fields...)
	case glog.Error:
		z.Error(msg, fields...)
	case glog.Disabled:
		// Nothing
	default:
		z.Error(msg, fields...)
	}
}

func (z adapter) Debug(msg string, fields ...any) {
	if !z.IsDebug() {
		return
	}

	z.lg.Debug(msg, toZap("", fields)...)
}

func (z adapter) Info(msg string, fields ...any) {
	if !z.IsInfo() {
		return
	}

	z.lg.Info(msg, toZap("", fields)...)
}

func (z adapter) Warn(msg string, fields ...any) {
	if !z.IsWarn() {
		return
	}

	z.lg.Warn(msg, toZap("", fields)...)
}

func (z adapter) Error(msg string, fields ...any) {
	if !z.IsError() {
		return
	}

	z.lg.Error(msg, toZap("", fields)...)
}

func (z adapter) IsDebug() bool {
	return z.max >= glog.Debug
}

func (z adapter) IsInfo() bool {
	return z.max >= glog.Info
}

func (z adapter) IsWarn() bool {
	return z.max >= glog.Warn
}

func (z adapter) IsError() bool {
	return z.max >= glog.Error
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
