package zcache

import (
	"context"

	"go.opentelemetry.io/otel/attribute"

	"github.com/hkoosha/giraffe/core/container/otel"
	"github.com/hkoosha/giraffe/core/t11y"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/glog"
)

const (
	Hit Outcome = iota + 1
	Miss
	Ignore
	BadKey
	BadValue
	BadData
	Bad
)

const (
	Get Op = iota + 1
	Set
	Unset
)

// ============================================================================.

type Outcome int

func (o Outcome) String() string {
	switch o {
	case Hit:
		return "hit"
	case Miss:
		return "miss"
	case Ignore:
		return "ignore"
	case BadKey:
		return "bad_key"
	case BadValue:
		return "bad_value"
	case BadData:
		return "bad_data"
	case Bad:
		return "unknown_error"
	default:
		panic(EF("unreachable"))
	}
}

// ============================================================================.

type Op int

func (o Op) String() string {
	switch o {
	case Get:
		return "get"
	case Set:
		return "set"
	case Unset:
		return "unset"
	default:
		panic(EF("unreachable"))
	}
}

// ============================================================================.

type Adapter[K comparable, V any] interface {
	Get(context.Context, K) (*Item[K, V], Outcome, error)

	Set(context.Context, K, V) (Outcome, error)

	Unset(context.Context, K) (Outcome, error)
}

type Item[K comparable, V any] struct {
	Key   K
	Value V
}

// ============================================================================.

func Of[K comparable, V any](
	adapter Adapter[K, V],
) *Cache[K, V] {
	return &Cache[K, V]{
		adapter:        adapter,
		lg:             nil,
		cnt:            nil,
		attrOtel:       nil,
		byOpAndOutcome: nil,
	}
}

type Cache[K comparable, V any] struct {
	lg             glog.Lg
	cnt            otel.Int64Counter
	adapter        Adapter[K, V]
	byOpAndOutcome map[Op]map[Outcome][]attribute.KeyValue
	attrOtel       []attribute.KeyValue
}

func (c *Cache[K, V]) WithLg(
	lg glog.Lg,
) *Cache[K, V] {
	t11y.NonNil(lg)
	cp := *&c
	cp.lg = lg
	return cp
}

func (c *Cache[K, V]) WithoutLg() *Cache[K, V] {
	cp := *&c
	cp.lg = nil
	return cp
}

func (c *Cache[K, V]) WithOtel(
	cnt otel.Int64Counter,
	attrs ...attribute.KeyValue,
) *Cache[K, V] {
	t11y.NonNil(cnt)

	cp := *&c

	cp.cnt = cnt
	cp.attrOtel = attrs

	hit := attribute.String("result", "hit")
	miss := attribute.String("result", "miss")
	ignore := attribute.String("result", "ignore")
	get := attribute.String("op", "get")
	set := attribute.String("op", "set")
	unset := attribute.String("op", "unset")

	cp.byOpAndOutcome = map[Op]map[Outcome][]attribute.KeyValue{
		Get: {
			Hit:    append([]attribute.KeyValue{get, hit}, attrs...),
			Miss:   append([]attribute.KeyValue{get, miss}, attrs...),
			Ignore: append([]attribute.KeyValue{get, ignore}, attrs...),
		},
		Set: {
			Hit:    append([]attribute.KeyValue{set, hit}, attrs...),
			Miss:   append([]attribute.KeyValue{set, miss}, attrs...),
			Ignore: append([]attribute.KeyValue{set, ignore}, attrs...),
		},
		Unset: {
			Hit:    append([]attribute.KeyValue{unset, hit}, attrs...),
			Miss:   append([]attribute.KeyValue{unset, miss}, attrs...),
			Ignore: append([]attribute.KeyValue{unset, ignore}, attrs...),
		},
	}

	return cp
}

func (c *Cache[K, V]) WithoutOtel() *Cache[K, V] {
	cp := *&c
	cp.cnt = nil
	cp.attrOtel = nil
	cp.byOpAndOutcome = nil
	return cp
}

func (c *Cache[K, V]) mkAttrs(
	op Op,
	outcome Outcome,
	err error,
) []attribute.KeyValue {
	//nolint:nestif
	if err == nil {
		if forOp, ok0 := c.byOpAndOutcome[op]; ok0 {
			if forOutcome, ok1 := forOp[outcome]; ok1 {
				return forOutcome
			}
		}
	}

	l := len(c.attrOtel)
	attr := make([]attribute.KeyValue, l+2)
	copy(attr, c.attrOtel)
	attr[l] = attribute.String("result", outcome.String())
	attr[l+1] = attribute.String("op", op.String())
	return attr
}

func (c *Cache[K, V]) Get(
	ctx context.Context,
	k K,
) *Item[K, V] {
	const op = Get

	item, result, err := c.adapter.Get(ctx, k)

	c.cnt.Inc(ctx, c.mkAttrs(op, result, err)...)

	if c.cnt != nil {
		c.cnt.Inc(ctx, c.mkAttrs(Set, result, err)...)
	}

	if err != nil && c.lg != nil {
		c.lg.Error(
			"cache error",
			N("op", op),
			N("result", result),
			N("key", k),
			err,
		)
	}

	if err != nil || result != Hit {
		return nil
	}
	return item
}

func (c *Cache[K, V]) Set(
	ctx context.Context,
	k K,
	v V,
) {
	const op = Set

	result, err := c.adapter.Set(ctx, k, v)

	if err != nil && c.lg != nil {
		c.lg.Error(
			"cache error",
			N("op", op),
			N("result", result),
			N("key", k),
			err,
		)
	}

	if c.cnt != nil {
		c.cnt.Inc(ctx, c.mkAttrs(Set, result, err)...)
	}
}

func (c *Cache[K, V]) Unset(
	ctx context.Context,
	k K,
) {
	const op = Unset

	result, err := c.adapter.Unset(ctx, k)

	if err != nil && c.lg != nil {
		c.lg.Error(
			"cache error",
			N("op", op),
			N("result", result),
			N("key", k),
			err,
		)
	}

	if c.cnt != nil {
		c.cnt.Inc(ctx, c.mkAttrs(op, result, err)...)
	}
}

// ============================================================================.

func Noop[K comparable, V any]() Adapter[K, V] {
	return &noop[K, V]{}
}
