package converters

import (
	"reflect"

	"github.com/hkoosha/giraffe/core/serdes/internal"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

var tConv = reflect.TypeOf(Bytes())

func IsConv(v any) bool {
	if v == nil {
		return false
	}

	ok, err := internal.ImplementsGenericErased(
		reflect.TypeOf(v),
		tConv,
	)

	return err == nil && ok
}

func Cast[T, U any](v any) (Conv[T, U], bool) {
	cast, ok := v.(Conv[T, U])
	return cast, ok
}

func MustCast[T, U any](v any) Conv[T, U] {
	cast, ok := v.(Conv[T, U])
	if !ok {
		var t T
		var u U
		panic(EF("not a converter, T=%v U=%v, v=%v", t, u, v))
	}
	return cast
}
