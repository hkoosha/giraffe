package zcache

import (
	"context"
)

type CacheClearedError struct{}

func (*CacheClearedError) Error() string { return "cache cleared" }

var ErrClearedCache = &CacheClearedError{}

// =============================================================================.

type Cache[K comparable, V any] interface {
	Get(context.Context, K) *Item[K, V]

	Set(context.Context, K, V)

	Delete(context.Context, K) error
}

type Item[K comparable, V any] struct {
	Key   K
	Value V
}
