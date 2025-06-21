package toggles

import (
	"context"
)

type Op uint

const (
	Eq Op = iota + 1
	NotEq
	In
	NotIn
	Has
	NotHas
)

//goland:noinspection GoUnusedExportedFunction
func Of[I interface {
	string | int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}](
	name string,
	op Op,
	value I,
) Attr {
	return newAttr(name, op, int64(value))
}

type Toggler interface {
	Get(
		ctx context.Context,
		name string,
		attrs ...Attr,
	) (bool, error)

	GetOrFalse(
		ctx context.Context,
		name string,
		attrs ...Attr,
	) bool

	Set(
		ctx context.Context,
		name string,
		enabled bool,
		attrs ...Attr,
	) error

	Enable(
		ctx context.Context,
		name string,
		attrs ...Attr,
	) error

	Disable(
		ctx context.Context,
		name string,
		attrs ...Attr,
	) error
}
