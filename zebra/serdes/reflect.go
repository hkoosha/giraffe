package serdes

import (
	"reflect"

	"github.com/hkoosha/giraffe/internal/reflected"
	. "github.com/hkoosha/giraffe/t11y/dot"
)

var tSerde = reflect.TypeOf(Bytes())

func IsSerde(v any) bool {
	if v == nil {
		return false
	}

	ok, err := reflected.ImplementsGenericErased(
		reflect.TypeOf(v),
		tSerde,
	)

	return err == nil && ok
}

func Cast[T any](v any) (Serde[T], bool) {
	cast, ok := v.(Serde[T])
	return cast, ok
}

func MustCast[T any](v any) Serde[T] {
	cast, ok := v.(Serde[T])
	if !ok {
		var t T
		panic(EF("not a serde, T=%v v=%v", t, v))
	}
	return cast
}
