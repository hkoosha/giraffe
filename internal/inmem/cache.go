package inmem

import (
	"math"
	"reflect"
	"sync"
	"time"

	"github.com/hkoosha/giraffe/internal/vendored/tlru"
)

const ttl = 4 * 24 * time.Hour

var cachesMu = &sync.Mutex{}
var caches = make(map[reflect.Type]any, 1)

type Item[V any] struct {
	Err error
	V   V
}

func (i Item[V]) Unpack() (V, error) {
	return i.V, i.Err
}

func getCache[V any]() *tlru.Cache[string, Item[V]] {
	var v V
	typ := reflect.TypeOf(v)

	cachesMu.Lock()
	defer cachesMu.Unlock()

	cache, ok := caches[typ]
	if !ok {
		cache = tlru.New[string, Item[V]](
			tlru.ConstantCost,
			math.MaxInt,
		)
		caches[typ] = cache
	}

	return cache.(*tlru.Cache[string, Item[V]])
}

func Get[V any](
	key string,
) (Item[V], bool) {
	cache := getCache[V]()

	if cached, _, ok := cache.Get(key); ok {
		return cached, true
	}

	return Item[V]{}, false
}

func Set[V any](
	key string,
	v V,
	err error,
) {
	item := Item[V]{
		V:   v,
		Err: err,
	}

	cache := getCache[V]()
	cache.Set(key, item, ttl)
}
