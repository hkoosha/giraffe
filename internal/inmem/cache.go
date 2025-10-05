package inmem

import (
	"math"
	"reflect"
	"sync"
	"time"

	"github.com/hkoosha/giraffe/internal/vendored/tlru"
	. "github.com/hkoosha/giraffe/t11y/dot"
)

const (
	ttl = 4 * 24 * time.Hour

	BucketParseQuery        = "parse_query"
	BucketReflectImplements = "reflect_implements"
)

var (
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

func Get[V any](
	bucket string,
	key string,
) (Item[V], bool) {
	cache := getCache[V](bucket)

	if cached, _, ok := cache.Get(key); ok {
		return cached, true
	}

	//nolint:exhaustruct
	return Item[V]{}, false
}

func Set[V any](
	bucket string,
	key string,
	v V,
	err error,
) {
	item := Item[V]{
		V:   v,
		Err: err,
	}

	cache := getCache[V](bucket)
	cache.Set(key, item, ttl)
}

func GetOr[V any](
	bucket string,
	key string,
	fn func() (V, error),
) (Item[V], error) {
	cached, ok := Get[V](bucket, key)

	var err error
	if !ok {
		v, fnErr := fn()
		err = fnErr
		Set(bucket, key, v, err)
		cached, ok = Get[V](bucket, key)
	}

	Assert(ok)

	if err != nil {
		return Item[V]{}, err
	}

	return cached, nil
}
