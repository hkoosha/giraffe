package finalizers

import (
	"context"
	"sync"

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/g11y/setup"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

type FinalizerFn00 = func()
type FinalizerFn10 = func(context.Context)
type FinalizerFn11 = func(context.Context) context.Context
type FinalizerFn interface {
	FinalizerFn11 | FinalizerFn00 | FinalizerFn10
}

func Add[F FinalizerFn](
	fn F,
) {
	AddTo(globalFin, fn)
}

func Finalize(
	ctx context.Context,
) context.Context {
	return globalFin.Finalize(ctx)
}

// ============================================================================.

func NewFinalizer(
	onceReg setup.Then,
) *Finalizer {
	g11y.NonNil(onceReg)

	return &Finalizer{
		onceReg:    onceReg,
		finalizers: []FinalizerFn11{},
		mu:         &sync.RWMutex{},
	}
}

func AddTo[F FinalizerFn](
	f *Finalizer,
	fn F,
) {
	switch v := any(fn).(type) {
	case FinalizerFn11:
		f.Add(v)

	case FinalizerFn10:
		f.Add10(v)

	case FinalizerFn00:
		f.Add00(v)

	default:
		panic(EF("unsupported finalizer type: %T", v))
	}
}

type Finalizer struct {
	mu         *sync.RWMutex
	finalizers []FinalizerFn11
	onceReg    setup.Then
}

func (f *Finalizer) Add00(
	fin FinalizerFn00,
) {
	f.Add(func(_ context.Context) context.Context {
		fin()
		return nil
	})
}

func (f *Finalizer) Add10(
	fin FinalizerFn10,
) {
	f.Add(func(ctx context.Context) context.Context {
		fin(ctx)
		return nil
	})
}

func (f *Finalizer) Add(
	fn FinalizerFn11,
) {
	g11y.NonNil(fn)
	f.add(fn)
}

func (f *Finalizer) Finalize(
	ctx context.Context,
) context.Context {
	fin := f.get()

	for i := len(fin) - 1; i >= 0; i-- {
		fn := fin[i]
		if ctx0 := fn(ctx); ctx0 != nil {
			ctx = ctx0
		}
	}

	return ctx
}
