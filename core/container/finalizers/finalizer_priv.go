package finalizers

import (
	"github.com/hkoosha/giraffe/core/container/setup"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

var globalFin = NewFinalizer(setup.Global())

func (f *Finalizer) ensure() {
	if f.onceReg == nil {
		panic(EF("finalizers in invalid state"))
	}
}

func (f *Finalizer) add(
	fn FinalizerFn1,
) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.ensure()
	f.finalizers = append(f.finalizers, fn)
}

func (f *Finalizer) get() []FinalizerFn1 {
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
