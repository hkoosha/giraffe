package zsync

import (
	"sync/atomic"
)

func valueOf[T any](v T) *atomic.Value {
	var av atomic.Value
	av.Store(v)
	return &av
}

func Of[T any](v T) *Value[T] {
	return &Value[T]{
		value: valueOf(v),
	}
}

type Value[T any] struct {
	value *atomic.Value
}

func (v *Value[T]) Load() T {
	//nolint:errcheck,forcetypeassert
	return v.value.Load().(T)
}

func (v *Value[T]) Store(t T) {
	v.value.Store(t)
}

//nolint:nonamedreturns
func (v *Value[T]) Swap(t T) (old T) {
	//nolint:errcheck,forcetypeassert
	return v.value.Swap(t).(T)
}

//nolint:nonamedreturns
func (v *Value[comparable]) CompareAndSwap(
	old any,
	newer any,
) (swapped bool) {
	return v.value.CompareAndSwap(old, newer)
}
