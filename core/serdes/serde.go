package serdes

import (
	"io"

	"github.com/hkoosha/giraffe/core/serdes/converters"
)

type Serde[T any] interface {
	converters.Conv[T, []byte]

	StreamTo(io.Writer, T) error

	StreamFrom(io.Reader) (T, error)
}

// =============================================================================.

func String() Serde[string] {
	return stringSerde{
		b: bytesSerde{},
	}
}

func Bytes() Serde[[]byte] {
	return bytesSerde{}
}

func Json[T any]() Serde[T] {
	return jsonSerde[T]{}
}
