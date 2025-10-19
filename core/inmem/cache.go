package inmem

import (
	"math"
	"reflect"
	"sync"
	"time"

	"github.com/hkoosha/giraffe/internal/vendored/tlru"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

var (
	bucketMu = &sync.Mutex{}
	buckets  = map[string]struct{}{}

	cachesMu = &sync.Mutex{}
	caches   = make(map[string]any, 1)
)

type Item[V any] struct {
	Err error
	V   V
}

func (i Item[V]) Unpack() (V, error) {
	return i.V, i.Err
}

func getCache[V any](
	bucket string,
) *tlru.Cache[string, Item[V]] {
	cachesMu.Lock()
	defer cachesMu.Unlock()

	cache, ok := caches[bucket]
	if !ok {
		cache = tlru.New[string, Item[V]](
			tlru.ConstantCost,
			math.MaxInt,
		)
		caches[bucket] = cache
	}

	cast, ok := cache.(*tlru.Cache[string, Item[V]])

	if !ok {
		var expecting *tlru.Cache[string, Item[V]]
		panic(EF(
			"unreachable: wrong data type, expecting=*%s got=%s",
			reflect.TypeOf(expecting).Elem().Name(),
			reflect.TypeOf(cache).Name(),
		))
	}

	return cast
}

type Cache[V any] struct {
	bucket string
	ttl    time.Duration
}

func (c Cache[V]) Get(
	key string,
) (Item[V], bool) {
	cache := getCache[V](c.bucket)

	if cached, _, ok := cache.Get(key); ok {
		return cached, true
	}

	//nolint:exhaustruct
	return Item[V]{}, false
}

func (c Cache[V]) Set(
	key string,
	v V,
	err error,
) {
	item := Item[V]{
		V:   v,
		Err: err,
	}

	cache := getCache[V](c.bucket)
	cache.Set(key, item, c.ttl)
}

func (c Cache[V]) GetOr(
	key string,
	fn func() (V, error),
) (Item[V], error) {
	cached, ok := c.Get(key)

	var err error
	if !ok {
		v, fnErr := fn()
		err = fnErr
		c.Set(key, v, err)
		cached, ok = c.Get(key)
	}

	Assert(ok)

	if err != nil {
		return Item[V]{}, err
	}

	return cached, nil
}

func Make[V any](
	bucket string,
	ttl time.Duration,
) Cache[V] {
	bucketMu.Lock()
	defer bucketMu.Unlock()

	if _, ok := buckets[bucket]; ok {
		panic(EF("cache already defined: %s", bucket))
	}
	buckets[bucket] = struct{}{}

	return Cache[V]{
		bucket: bucket,
		ttl:    ttl,
	}
}
