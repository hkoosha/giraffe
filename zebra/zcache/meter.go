package zcache

import (
	"context"

	"go.opentelemetry.io/otel/attribute"

	"github.com/hkoosha/giraffe/zebra/serdes"
)

const (
	o11yErrKey   = "err"
	o11yOkPrefix = "ok_"
)

var _ Cache[any, any] = (*metered[any, any])(nil)

type metered[K comparable, V any] struct {
	adapter   Adapter[K, V]
	keyConv   serdes.Conv[K, string]
	valConv   serdes.Conv[V, string]
	o11yAttrs []attribute.KeyValue
	cnt       *otel.Counter
}

func (r *metered[K, V]) keyOf(k K) string {
	key, err := r.keyConv.Write(k)
	if err != nil {
		return o11yErrKey
	}

	return o11yOkPrefix + key
}

func (r *metered[K, V]) Get(
	ctx context.Context,
	k K,
) *Item[K, V] {
	item, result, err := r.adapter.Get(ctx, k)
	if err != nil {
		return nil
	}

	return item
}

func (r *metered[K, V]) Set(
	ctx context.Context,
	k K,
	v V,
) {
	panic("todo")
}

func (r *metered[K, V]) Delete(
	ctx context.Context,
	k K,
) error {
	panic("todo")
}
