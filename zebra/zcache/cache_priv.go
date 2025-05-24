package zcache

import (
	"context"
)

type noop[K comparable, V any] struct{}

func (n noop[K, V]) Get(
	context.Context,
	K,
) (*Item[K, V], Outcome, error) {
	return nil, Miss, nil
}

func (n noop[K, V]) Set(
	context.Context,
	K,
	V,
) (Outcome, error) {
	return Ignore, nil
}

func (n noop[K, V]) Unset(
	context.Context,
	K,
) (Outcome, error) {
	return Ignore, nil
}
