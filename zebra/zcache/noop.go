package zcache

import "context"

func Noop[K comparable, V any]() Cache[K, V] {
	return &noopCache[K, V]{}
}

type noopCache[K comparable, V any] struct{}

func (n noopCache[K, V]) Get(context.Context, K) *Item[K, V] {
	return nil
}

func (n noopCache[K, V]) Set(context.Context, K, V) {
}

func (n noopCache[K, V]) Delete(context.Context, K) error {
	return nil
}
