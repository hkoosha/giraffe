package internal

import (
	"reflect"
)

var (
	TISize = reflect.TypeOf((*int)(nil)).Elem()
	TI8    = reflect.TypeOf((*int8)(nil)).Elem()
	TI16   = reflect.TypeOf((*int16)(nil)).Elem()
	TI32   = reflect.TypeOf((*int32)(nil)).Elem()
	TI64   = reflect.TypeOf((*int64)(nil)).Elem()

	TUSize = reflect.TypeOf((*uint)(nil)).Elem()
	TU8    = reflect.TypeOf((*uint8)(nil)).Elem()
	TU16   = reflect.TypeOf((*uint16)(nil)).Elem()
	TU32   = reflect.TypeOf((*uint32)(nil)).Elem()
	TU64   = reflect.TypeOf((*uint64)(nil)).Elem()

	TStr = reflect.TypeOf((*string)(nil)).Elem()
)
