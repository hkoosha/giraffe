package zcache

import (
	"context"
)

const (
	CacheOpSuccess CacheOpResult = iota + 1
	CacheOpMiss
	CacheOpHit
	CacheOpBadKey
	CacheOpBadValue
	CacheOpBadData
)

type CacheOpResult int

type Adapter[K comparable, V any] interface {
	Get(context.Context, K) (*Item[K, V], CacheOpResult, error)

	Set(context.Context, K, V) (CacheOpResult, error)

	Delete(context.Context, K) (CacheOpResult, error)
}

type Cache[K comparable, V any] interface {
	Get(context.Context, K) *Item[K, V]

	Set(context.Context, K, V)

	Delete(context.Context, K) error
}

type Item[K comparable, V any] struct {
	Key   K
	Value V
}
