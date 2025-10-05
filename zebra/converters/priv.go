package converters

import (
	"encoding/json"
)

type bytesConv struct{}

func (s bytesConv) Write(b []byte) ([]byte, error) {
	return b, nil
}

func (s bytesConv) Read(b []byte) ([]byte, error) {
	return b, nil
}

// =============================================================================

type jsonConv[T any] struct{}

func (s jsonConv[T]) Write(t T) ([]byte, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s jsonConv[T]) Read(b []byte) (T, error) {
	var t T
	err := json.Unmarshal(b, &t)
	return t, err
}

// =============================================================================

type jsonStr[T any] struct{}

func (s jsonStr[T]) Write(t T) (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s jsonStr[T]) Read(b string) (T, error) {
	var t T
	err := json.Unmarshal([]byte(b), &t)
	return t, err
}

// =============================================================================

type stringConv struct{}

func (s stringConv) Write(t string) (string, error) {
	return t, nil
}

func (s stringConv) Read(t string) (string, error) {
	return t, nil
}
