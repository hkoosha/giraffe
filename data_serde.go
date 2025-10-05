package giraffe

import (
	"bytes"
	"encoding/json"
	"io"

	. "github.com/hkoosha/giraffe/t11y/dot"
	"github.com/hkoosha/giraffe/zebra/serdes"
)

type datumSerde struct{}

func (s datumSerde) Write(v Datum) ([]byte, error) {
	return v.MarshalJSON()
}

func (s datumSerde) Read(b []byte) (Datum, error) {
	var v any
	if err := json.Unmarshal(b, v); err != nil {
		return OfErr(), E(err)
	}

	return ofJsonable(v)
}

func (s datumSerde) StreamTo(w io.Writer, v Datum) error {
	enc := json.NewEncoder(w)
	return E(enc.Encode(v))
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
