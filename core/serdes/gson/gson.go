package gson

import (
	"bytes"
	"encoding/json"
	"io"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

func MustMarshal(
	v any,
) []byte {
	return M(Marshal(v))
}

func MustUnmarshal[V any](
	b []byte,
) V {
	return M(Unmarshal[V](b))
}

func Marshal(
	v any,
) ([]byte, error) {
	return T(json.Marshal(v))
}

func Unmarshal[V any](
	b []byte,
) (V, error) {
	var v V

	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()

	err := E(decoder.Decode(&v))
	return v, err
}

func NewEncoder(w io.Writer) *json.Encoder {
	enc := json.NewEncoder(w)
	return enc
}

func EncodeTo(
	w io.Writer,
	v any,
) error {
	return E(NewEncoder(w).Encode(v))
}
