package setup

import (
	"context"
	"sync"

	. "github.com/hkoosha/giraffe/dot"
)

type Finalizer = func(context.Context)

func NewFinalizerRegistry(
	name string,
) *FinalizerRegistry {
	Once("boot", "setup", "finalizer", name)

	return &FinalizerRegistry{
		name:       name,
		finalizers: []Finalizer{},
		mu:         &sync.Mutex{},
	}
}

type FinalizerRegistry struct {
	mu         *sync.Mutex
	name       string
	finalizers []Finalizer
}

func (f *FinalizerRegistry) Add(
	fin Finalizer,
) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.finalizers == nil {
		panic(EF("finalizers already executed"))
	}

	f.finalizers = append(f.finalizers, fin)
}

func (f *FinalizerRegistry) Execute(
	ctx context.Context,
) {
	// Prevent deadlocks on mutex reentry.
	ch := make(chan []Finalizer, 1)
	go func(chan<- []Finalizer) {
		f.mu.Lock()
		defer f.mu.Unlock()

		fn := f.finalizers
		f.finalizers = nil
		ch <- fn
	}(ch)

	fin := <-ch

	if fin == nil {
		panic(EF("finalizers already executed"))
	}

	for _, fn := range fin {
		fn(ctx)
	}
}
