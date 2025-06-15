package gzapadapter

import (
	"strconv"

	"go.uber.org/zap"

	"github.com/hkoosha/giraffe/glog"
)

func Of(lg *zap.Logger) glog.GLog {
	return &adapter{lg: lg}
}

// ============================================================================.

type adapter struct {
	lg *zap.Logger
}

func (z adapter) Named(s string) glog.GLog {
	return adapter{lg: z.lg.Named(s)}
}

func (z adapter) Debug(msg string, fields ...any) {
	z.lg.Debug(msg, toZap(fields)...)
}

func (z adapter) Info(msg string, fields ...any) {
	z.lg.Info(msg, toZap(fields)...)
}

func (z adapter) Warn(msg string, fields ...any) {
	z.lg.Warn(msg, toZap(fields)...)
}

func (z adapter) Error(msg string, fields ...any) {
	z.lg.Error(msg, toZap(fields)...)
}

func toZap(fields []any) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		switch v := f.(type) {
		case zap.Field:
			zapFields[i] = v
		default:
			zapFields[i] = zap.Any("field"+strconv.Itoa(i), v)
		}
	}
	return zapFields
}
