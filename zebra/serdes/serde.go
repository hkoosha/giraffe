package serdes

import (
	"bytes"
	"encoding/json"
	"io"
)

type Conv[T any, U any] interface {
	Write(T) (U, error)
	Read(U) (T, error)
}

type Serde[T any] interface {
	Conv[T, []byte]

	WriteTo(T, io.Writer) error
	ReadFrom(io.Reader) (T, error)
}

// =============================================================================.

func Bytes() Serde[[]byte] {
	return &bytesSerde{}
}

type bytesSerde struct{}

func (s bytesSerde) Write(b []byte) ([]byte, error) {
	return b, nil
}

func (s bytesSerde) WriteTo(b []byte, w io.Writer) error {
	n, err := io.Copy(w, bytes.NewReader(b))
	if err != nil {
		return err
	}

	if n != int64(len(b)) {
		return errTruncatedStream
	}

	return nil
}

func (s bytesSerde) Read(b []byte) ([]byte, error) {
	return b, nil
}

//goland:noinspection GoStandardMethods
func (s bytesSerde) ReadFrom(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// =============================================================================.

func Json[T any]() Serde[T] {
	return &jsonSerde[T]{}
}

type jsonSerde[T any] struct{}

func (s jsonSerde[T]) Write(t T) ([]byte, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s jsonSerde[T]) WriteTo(t T, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}

func (s jsonSerde[T]) Read(b []byte) (T, error) {
	return s.ReadFrom(bytes.NewReader(b))
}

//goland:noinspection GoStandardMethods
func (s jsonSerde[T]) ReadFrom(r io.Reader) (T, error) {
	var t T
	if err := json.NewDecoder(r).Decode(&t); err != nil {
		return t, err
	}

	return t, nil
}

// =============================================================================.

func JsonConv[T any]() Conv[T, string] {
	return &jsonStr[T]{jsonSerde[T]{}}
}

type jsonStr[T any] struct {
	jsonSerde[T]
}

func (s *jsonStr[T]) Write(t T) (string, error) {
	b, err := s.jsonSerde.Write(t)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (s *jsonStr[T]) Read(js string) (T, error) {
	return s.jsonSerde.Read([]byte(js))
}

// =====================================.

func String() Serde[string] {
	return &stringSerde{}
}

type stringSerde struct{}

func (s stringSerde) Write(t string) ([]byte, error) {
	return []byte(t), nil
}

func (s stringSerde) WriteTo(t string, w io.Writer) error {
	n, err := io.WriteString(w, t)
	if err != nil {
		return err
	}

	if n != len(t) {
		return errTruncatedStream
	}

	return nil
}

func (s stringSerde) Read(b []byte) (string, error) {
	return string(b), nil
}

//goland:noinspection GoStandardMethods
func (s stringSerde) ReadFrom(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// =====================================.

func StringConv() Conv[string, string] {
	return &stringConv{}
}

type stringConv struct{}

func (s stringConv) Write(t string) (string, error) {
	return t, nil
}

func (s stringConv) Read(t string) (string, error) {
	return t, nil
}
