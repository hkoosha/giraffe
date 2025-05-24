package serdes

import (
	"reflect"

	"github.com/hkoosha/giraffe/core/serdes/internal"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

var tSerde = reflect.TypeOf(Bytes())

func IsSerde(v any) bool {
	if v == nil {
		return false
	}

	ok, err := internal.ImplementsGenericErased(
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
		panic(EF(
			"not a serde or incompatible type, T=%s v=%s",
			reflect.TypeOf(t).String(),
			reflect.TypeOf(v).String(),
		))
	}
	return cast
}
