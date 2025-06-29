package zcache

import (
	"github.com/hkoosha/giraffe/g11y/gtx"
)

type noop[K comparable, V any] struct{}

func (n noop[K, V]) Get(
	gtx.Context,
	K,
) (*Item[K, V], Outcome, error) {
	return nil, Miss, nil
}

func (n noop[K, V]) Set(
	gtx.Context,
	K,
	V,
) (Outcome, error) {
	return Ignore, nil
}

func (n noop[K, V]) Unset(
	gtx.Context,
	K,
) (Outcome, error) {
	return Ignore, nil
}
