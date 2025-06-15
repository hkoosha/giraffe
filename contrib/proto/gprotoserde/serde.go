package gprotoserde

import (
	"io"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/zebra/serdes"
)

func New[T proto.Message]() serdes.Serde[T] {
	return Of[T](
		protojson.UnmarshalOptions{},
		protojson.MarshalOptions{},
	)
}

func Of[T proto.Message](
	unmarshal protojson.UnmarshalOptions,
	marshal protojson.MarshalOptions,
) serdes.Serde[T] {
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

func (s serde[T]) WriteTo(t T, w io.Writer) error {
	b, err := s.Write(t)
	if err != nil {
		return err
	}

	n, err := w.Write(b)
	if err != nil {
		return err
	}

	if n != len(b) {
		panic(EF("short write: buffer_len=%d, written=%d", len(b), n))
	}

	return nil
}

func (s serde[T]) Read(b []byte) (T, error) {
	var t T
	err := s.unmarshal.Unmarshal(b, t)
	return t, err
}

//goland:noinspection GoStandardMethods
func (s serde[T]) ReadFrom(r io.Reader) (T, error) {
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
