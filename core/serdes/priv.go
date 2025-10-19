package serdes

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/hkoosha/giraffe/core/serdes/gson"
)

var errTruncatedStream = errors.New("truncated stream")

// =============================================================================

type stringSerde struct {
	b bytesSerde
}

func (s stringSerde) Write(t string) ([]byte, error) {
	return []byte(t), nil
}

func (s stringSerde) Read(b []byte) (string, error) {
	return string(b), nil
}

func (s stringSerde) StreamTo(w io.Writer, v string) error {
	return s.b.StreamTo(w, []byte(v))
}

func (s stringSerde) StreamFrom(r io.Reader) (string, error) {
	v, err := s.b.StreamFrom(r)
	if err != nil {
		return "", err
	}

	return string(v), nil
}

// =============================================================================

type bytesSerde struct{}

func (s bytesSerde) Write(v []byte) ([]byte, error) {
	return v, nil
}

func (s bytesSerde) Read(v []byte) ([]byte, error) {
	return v, nil
}

func (s bytesSerde) StreamTo(w io.Writer, v []byte) error {
	n, err := io.Copy(w, bytes.NewReader(v))
	if err != nil {
		return err
	}
	if n != int64(len(v)) {
		return errTruncatedStream
	}

	return nil
}

func (s bytesSerde) StreamFrom(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// =============================================================================

type jsonSerde[T any] struct{}

func (s jsonSerde[T]) Write(v T) ([]byte, error) {
	return gson.Marshal(v)
}

func (s jsonSerde[T]) Read(v []byte) (T, error) {
	return gson.Unmarshal[T](v)
}

func (s jsonSerde[T]) StreamTo(w io.Writer, v T) error {
	return json.NewEncoder(w).Encode(v)
}

func (s jsonSerde[T]) StreamFrom(r io.Reader) (T, error) {
	var v T
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return v, err
	}

	return v, nil
}
