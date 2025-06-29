package finalizers

import (
	"context"

	"github.com/hkoosha/giraffe/g11y/setup"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

var globalFin = NewFinalizer(setup.Global())

func (f *Finalizer) ensure() {
	if f.onceReg == nil {
		panic(EF("finalizers in invalid state"))
	}
}

func (f *Finalizer) add(
	fn func(context.Context) context.Context,
) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.ensure()
	f.finalizers = append(f.finalizers, fn)
}

func (f *Finalizer) get() []func(context.Context) context.Context {
	// Important: if the lock is moved outside of this function (e.g. in
	// Finalize()), it will can potentially cause a deadlock.
	f.mu.Lock()
	defer f.mu.Unlock()

	f.ensure()

	fin := f.finalizers
	f.finalizers = nil
	f.onceReg = nil

	return fin
}
