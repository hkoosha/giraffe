package gprotoserde

import (
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/hkoosha/giraffe/core/t11y"
)

type Conv[T any, U any] interface {
	Write(T) (U, error)
	Read(U) (T, error)
}

type Serde[T any] interface {
	Conv[T, []byte]

	StreamTo(io.Writer, T) error
	StreamFrom(io.Reader) (T, error)
}

func New[T proto.Message]() Serde[T] {
	return Of[T](
		protojson.UnmarshalOptions{},
		protojson.MarshalOptions{},
	)
}

func Of[T proto.Message](
	unmarshal protojson.UnmarshalOptions,
	marshal protojson.MarshalOptions,
) Serde[T] {
	return &serde[T]{
		unmarshal: unmarshal,
		marshal:   marshal,
	}
}

// ============================================================================.

type serde[T proto.Message] struct {
	unmarshal protojson.UnmarshalOptions
	marshal   protojson.MarshalOptions
}

func (s serde[T]) Write(t T) ([]byte, error) {
	return s.marshal.Marshal(t)
}

func (s serde[T]) StreamTo(w io.Writer, t T) error {
	b, err := s.Write(t)
	if err != nil {
		return err
	}

	n, err := w.Write(b)
	if err != nil {
		return err
	}

	if n != len(b) {
		panic(t11y.TracedFmt("short write: buffer_len=%d, written=%d", len(b), n))
	}

	return nil
}

func (s serde[T]) Read(b []byte) (T, error) {
	var t T
	err := s.unmarshal.Unmarshal(b, t)
	return t, err
}

func (s serde[T]) StreamFrom(r io.Reader) (T, error) {
	var t T

	b, err := io.ReadAll(r)
	if err != nil {
		return t, err
	}

	err = s.unmarshal.Unmarshal(b, t)
	if err != nil {
		return t, err
	}

	return t, nil
}
