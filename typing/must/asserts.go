package must

import (
	"reflect"
	"slices"

	"github.com/hkoosha/giraffe/internal/dot0"
)

func IsOneOf(
	v any,
	allowedValues ...any,
) {
	if !slices.Contains(allowedValues, v) {
		panic(dot0.EF(
			"invalid value: %s#%v, expecting one of: %v",
			reflect.TypeOf(v).String(),
			v,
			allowedValues,
		))
	}
}
