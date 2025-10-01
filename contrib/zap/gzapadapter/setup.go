package gzapadapter

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var writeSyncer = sync.OnceValue(func() zapcore.WriteSyncer {
	return zapcore.Lock(os.Stdout)
})

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

func DefaultShutdown() {
	if err := zap.L().Sync(); err != nil {
		fmt.Println("Failed to sync logger:", err)
	}
}

func DefaultNewLogger(
	json bool,
	cfg *zapcore.EncoderConfig,
	extra ...zap.Field,
) *zap.Logger {
	level := zap.InfoLevel

	var enc zapcore.Encoder
	if json {
		enc = zapcore.NewJSONEncoder(*cfg)
	} else {
		enc = zapcore.NewConsoleEncoder(*cfg)
	}

	core := zapcore.NewCore(enc, writeSyncer(), level)

	return zap.New(
		core,
		zap.WithCaller(false),
		zap.AddStacktrace(zap.DPanicLevel),
		zap.Fields(extra...),
	)
}

func DefaultSetup(
	json bool,
	extra ...zap.Field,
) *zap.Logger {
	zap.ReplaceGlobals(DefaultNewLogger(
		json,
		DefaultEncoderConfig(),
	))

	return zap.L()
}
