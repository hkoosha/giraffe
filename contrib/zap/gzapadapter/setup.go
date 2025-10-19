package gzapadapter

import (
	"fmt"
	"maps"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

var writeSyncer = sync.OnceValue(func() zapcore.WriteSyncer {
	return zapcore.Lock(os.Stdout)
})

type Provider[D any] interface {
	Close() error

	Open() (D, error)

	Get() D

	Init(
		local bool,
		adjustGlobalLogger bool,
		extra map[string]string,
	)
}

type provider struct {
	mu                 *sync.Mutex
	lg                 *zap.Logger
	extra              map[string]string
	local              bool
	adjustGlobalLogger bool
	initialized        bool
	ready              bool
}

func (p *provider) Init(
	local bool,
	adjustGlobalLogger bool,
	extra map[string]string,
) {
	extra = maps.Clone(extra)
	if extra == nil {
		extra = make(map[string]string)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.ready = false
	p.initialized = false
	p.local = local
	p.adjustGlobalLogger = adjustGlobalLogger
	p.extra = extra
	p.initialized = true
}

func (p *provider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.initialized = false
	p.ready = false

	err := p.lg.Sync()
	if err != nil {
		//nolint:forbidigo // is a last effort, ok to print directly as logger is failing anyway
		fmt.Println("failed to sync logger:", err)
	}

	return err
}

func (p *provider) Open() (*zap.Logger, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	const level = zap.InfoLevel

	if !p.initialized {
		panic(EF("not initialized"))
	}

	if p.ready {
		Assert(p.lg != nil)
		return p.lg, nil
	}

	//nolint:exhaustruct // defaults are good enough
	cfg := zapcore.EncoderConfig{
		MessageKey:     "what",
		NameKey:        "who",
		TimeKey:        "when",
		CallerKey:      "where",
		StacktraceKey:  "whereabouts",
		LevelKey:       "level",
		FunctionKey:    "fn",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var enc zapcore.Encoder
	if p.local {
		enc = zapcore.NewConsoleEncoder(cfg)
	} else {
		enc = zapcore.NewJSONEncoder(cfg)
	}

	op := []zap.Option{
		zap.WithCaller(false),
		zap.AddStacktrace(zap.DPanicLevel),
	}

	if len(p.extra) > 0 {
		var fields []zap.Field
		for k, v := range p.extra {
			fields = append(fields, zap.String(k, v))
		}
		op = append(op, zap.Fields(fields...))
	}

	p.lg = zap.New(
		zapcore.NewCore(enc, writeSyncer(), level),
		op...,
	)

	if p.adjustGlobalLogger {
		zap.ReplaceGlobals(p.lg)
	}

	p.ready = true
	return p.lg, nil
}

func (p *provider) Get() *zap.Logger {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.initialized {
		panic(EF("not initialized"))
	}
	if !p.ready {
		panic(EF("not opened"))
	}

	return p.lg
}

func MkProvider() Provider[*zap.Logger] {
	return &provider{
		mu:                 &sync.Mutex{},
		lg:                 nil,
		extra:              map[string]string{},
		local:              false,
		initialized:        false,
		adjustGlobalLogger: false,
		ready:              false,
	}
}

func MkInit(
	local bool,
	adjustGlobalLogger bool,
	extra map[string]string,
) Provider[*zap.Logger] {
	p := MkProvider()
	p.Init(local, adjustGlobalLogger, extra)
	return p
}
