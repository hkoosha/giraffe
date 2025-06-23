package gzapadapter

import (
	"go.uber.org/zap/zapcore"
)

func DefaultEncoderConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		TimeKey:        "at",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		FunctionKey:    "fn",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
