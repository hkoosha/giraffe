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
	Write(T) ([]byte, error)
	WriteTo(T, io.Writer) error

	Read([]byte) (T, error)
	ReadFrom(io.Reader) (T, error)
}

// =============================================================================.

func JsonSerde[T any]() Serde[T] {
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
	return &jsonConv[T]{}
}

type jsonConv[T any] struct {
	j jsonSerde[T]
}

func (s jsonConv[T]) Write(t T) (string, error) {
	b, err := s.j.Write(t)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (s jsonConv[T]) Read(js string) (T, error) {
	var t T
	if err := json.Unmarshal([]byte(js), &t); err != nil {
		return t, err
	}

	return t, nil
}

// =====================================.

func StringSerde() Serde[string] {
	return &strSerde{}
}

type strSerde struct{}

func (s strSerde) Write(t string) ([]byte, error) {
	return []byte(t), nil
}

func (s strSerde) WriteTo(t string, w io.Writer) error {
	n, err := io.WriteString(w, t)
	if err != nil {
		return err
	}

	if n != len(t) {
		panic("unreachable: not all of the string was written")
	}

	return nil
}

func (s strSerde) Read(b []byte) (string, error) {
	return string(b), nil
}

//goland:noinspection GoStandardMethods
func (s strSerde) ReadFrom(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// =====================================.

func StringConv() Conv[string, string] {
	return &strConv{}
}

type strConv struct{}

func (s strConv) Write(t string) (string, error) {
	return t, nil
}

func (s strConv) Read(t string) (string, error) {
	return t, nil
}
