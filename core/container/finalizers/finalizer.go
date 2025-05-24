package finalizers

import (
	"context"
	"sync"

	"github.com/hkoosha/giraffe/core/t11y"

	"github.com/hkoosha/giraffe/core/container/setup"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

type (
	FinalizerFn0 = func()
	FinalizerFn1 = func(context.Context)
	FinalizerFn  interface {
		FinalizerFn0 | FinalizerFn1
	}
)

func Add[F FinalizerFn](
	fn F,
) {
	AddTo(globalFin, fn)
}

func Finalize(
	ctx context.Context,
) {
	globalFin.Finalize(ctx)
}

// ============================================================================.

func NewFinalizer(
	onceReg setup.Registry,
) *Finalizer {
	t11y.NonNil(onceReg)

	return &Finalizer{
		onceReg:    onceReg,
		finalizers: []FinalizerFn1{},
		mu:         &sync.RWMutex{},
	}
}

func AddTo[F FinalizerFn](
	f *Finalizer,
	fn F,
) {
	switch v := any(fn).(type) {
	case FinalizerFn1:
		f.Add10(v)

	case FinalizerFn0:
		f.Add00(v)

	default:
		panic(EF("unsupported finalizer type: %T", v))
	}
}

type Finalizer struct {
	onceReg    setup.Registry
	mu         *sync.RWMutex
	finalizers []FinalizerFn1
}

func (f *Finalizer) Add00(
	fin FinalizerFn0,
) {
	f.Add(func(_ context.Context) {
		fin()
	})
}

func (f *Finalizer) Add10(
	fin FinalizerFn1,
) {
	f.Add(func(ctx context.Context) {
		fin(ctx)
	})
}

func (f *Finalizer) Add(
	fn FinalizerFn1,
) {
	t11y.NonNil(fn)
	f.add(fn)
}

func (f *Finalizer) Finalize(
	ctx context.Context,
) {
	fin := f.get()

	for i := len(fin) - 1; i >= 0; i-- {
		fn := fin[i]
		fn(ctx)
	}
}
