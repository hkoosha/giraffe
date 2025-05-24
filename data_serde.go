package giraffe

import (
	"bytes"
	"io"

	"github.com/hkoosha/giraffe/core/serdes"
	"github.com/hkoosha/giraffe/core/serdes/gson"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

type datumSerde struct{}

func (s datumSerde) Write(v Datum) ([]byte, error) {
	return v.MarshalJSON()
}

func (s datumSerde) Read(b []byte) (Datum, error) {
	v, err := gson.Unmarshal[any](b)
	if err != nil {
		return OfErr(), err
	}
	return ofJsonable(v)
}

func (s datumSerde) StreamTo(w io.Writer, v Datum) error {
	return gson.EncodeTo(w, v)
}

func (s datumSerde) StreamFrom(r io.Reader) (Datum, error) {
	// TODO optimize this horrible impl.

	payload := new(bytes.Buffer)
	if _, err := io.Copy(payload, r); err != nil {
		return OfErr(), E(err)
	}

	return s.Read(payload.Bytes())
}

func DatumSerde() serdes.Serde[Datum] {
	return datumSerde{}
}
