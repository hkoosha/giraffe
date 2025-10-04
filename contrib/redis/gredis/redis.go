package gredis

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"

	"github.com/hkoosha/giraffe/g11y/glog"
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/g11y/t11y"
	"github.com/hkoosha/giraffe/zebra/serdes"
	"github.com/hkoosha/giraffe/zebra/zcache"
)

func New[K comparable, V any](
	cfg *Config,
	keySerde serdes.Conv[K, string],
	valSerde serdes.Conv[V, string],
) zcache.Adapter[K, V] {
	t11y.NonNil(cfg, keySerde, valSerde)
	cfg.Ensure()

	return &adapter[K, V]{
		cfg:      cfg,
		keySerde: keySerde,
		valSerde: valSerde,
	}
}

func NewForStringK[V any](
	cfg *Config,
	valSerde serdes.Conv[V, string],
) zcache.Adapter[string, V] {
	return New(
		cfg,
		serdes.StringConv(),
		valSerde,
	)
}

func NewForString(
	cfg *Config,
) zcache.Adapter[string, string] {
	return New(
		cfg,
		serdes.StringConv(),
		serdes.StringConv(),
	)
}

func NewForJson[V any](
	cfg *Config,
) zcache.Adapter[string, V] {
	// TODO remove hard dependency.
	if reflect.TypeFor[V]().Implements(reflect.TypeFor[proto.Message]()) {
		panic(EF(
			"cannot use json serde for proto values, use proto serde instead",
		))
	}

	return New(
		cfg,
		serdes.StringConv(),
		serdes.JsonConv[V, string](),
	)
}

// =============================================================================.

var _ zcache.Adapter[string, any] = (*adapter[string, any])(nil)

type adapter[K comparable, V any] struct {
	lg       glog.Lg
	cfg      *Config
	keySerde serdes.Conv[K, string]
	valSerde serdes.Conv[V, string]
	rds      *redis.Client
}

func (r *adapter[K, V]) keyOf(k K) (string, error) {
	key, err := r.keySerde.Write(k)
	if err != nil {
		return "", err
	}

	if r.cfg.namespace != "" {
		key = r.cfg.namespace + key
	}

	return key, nil
}

func (r *adapter[K, V]) Get(
	ctx context.Context,
	k K,
) (*zcache.Item[K, V], zcache.Outcome, error) {
	key, err := r.keyOf(k)
	if err != nil {
		return nil, zcache.BadKey, err
	}

	ctx, cancel := r.start(ctx)
	defer cancel()
	cmd := r.rds.Get(ctx, key)

	switch {
	case cmd == nil,
		errors.Is(cmd.Err(), redis.Nil):
		return nil, zcache.Miss, nil

	case cmd.Err() != nil:
		return nil, zcache.Bad, cmd.Err()
	}

	val, err := r.valSerde.Read(cmd.Val())
	if err != nil {
		return nil, zcache.BadData, err
	}

	return &zcache.Item[K, V]{Key: k, Value: val}, zcache.Hit, nil
}

func (r *adapter[K, V]) Set(
	ctx context.Context,
	k K,
	v V,
) (zcache.Outcome, error) {
	key, err := r.keyOf(k)
	if err != nil {
		return zcache.BadKey, err
	}

	val, err := r.valSerde.Write(v)
	if err != nil {
		return zcache.BadValue, err
	}

	ctx, cancel := r.start(ctx)
	defer cancel()
	cmd := r.rds.Set(ctx, key, val, r.cfg.TTL())

	if cmd != nil && cmd.Err() != nil {
		return zcache.Bad, cmd.Err()
	}

	return zcache.Hit, nil
}

func (r *adapter[K, V]) Unset(
	ctx context.Context,
	k K,
) (zcache.Outcome, error) {
	key, err := r.keyOf(k)
	if err != nil {
		return zcache.BadKey, err
	}

	ctx, cancel := r.start(ctx)
	defer cancel()
	cmd := r.rds.Del(ctx, key)

	if cmd != nil && cmd.Err() != nil {
		return zcache.Bad, cmd.Err()
	}

	return zcache.Hit, nil
}

func (r *adapter[K, V]) start(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	return context.WithDeadline(ctx, time.Now().Add(r.cfg.timeout))
}
