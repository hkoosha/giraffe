package serdes

import (
	"encoding/json"
)

type Serde[T any, U any] interface {
	Write(t T) U
	Read(s U) (T, error)
}

// =============================================================================.

func JSONSerde[T any]() Serde[T, string] {
	return &jsonSerde[T]{}
}

type jsonSerde[T any] struct{}

func (s jsonSerde[T]) Write(t T) string {
	data, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func (s jsonSerde[T]) Read(js string) (T, error) {
	var t T
	if err := json.Unmarshal([]byte(js), &t); err != nil {
		return t, err
	}

	return t, nil
}

// =====================================.

func StringSerde() Serde[string, string] {
	return &stringSerde{}
}

type stringSerde struct{}

func (s stringSerde) Write(t string) string {
	return t
}

func (s stringSerde) Read(t string) (string, error) {
	return t, nil
}
