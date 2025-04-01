package typing

import (
	"reflect"
)

//nolint:exhaustive
func IsNillable(kind reflect.Kind) bool {
	switch kind {
	case reflect.Ptr,
		reflect.Map,
		reflect.Chan,
		reflect.Func,
		reflect.Slice,
		reflect.Interface:
		return true

	default:
		return false
	}
}
