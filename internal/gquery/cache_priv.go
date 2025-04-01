package gquery

import (
	"math"
	"time"

	"github.com/hkoosha/giraffe/internal/vendored/tlru"
)

var (
	ttl = 1 * time.Hour

	cache = tlru.New[string, QueryCacheItem](
		tlru.ConstantCost,
		math.MaxInt,
	)

	noCache = QueryCacheItem{
		Query: ErrQ,
		Error: nil,
	}
)

type QueryCacheItem struct {
	Error error
	Query Query
}

func (i QueryCacheItem) Unpack() (Query, error) {
	return i.Query, i.Error
}

func get(
	spec string,
) (QueryCacheItem, bool) {
	if cached, _, ok := cache.Get(spec); ok {
		return cached, true
	}

	return noCache, false
}

func set(
	spec string,
	query Query,
	err error,
) QueryCacheItem {
	item := QueryCacheItem{
		Query: query,
		Error: err,
	}

	cache.Set(spec, item, ttl)

	return item
}
