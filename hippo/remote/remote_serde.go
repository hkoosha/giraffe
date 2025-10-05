package remote

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/gson"
	. "github.com/hkoosha/giraffe/t11y/dot"
	"github.com/hkoosha/giraffe/zebra/serdes"
)

type requestCompensations struct {
	With     any     `json:"with"                 yaml:"with"`
	OnErrRe  *string `json:"on_err_re,omitempty"  yaml:"on_err_re,omitempty"`
	OnNameRe *string `json:"on_name_re,omitempty" yaml:"on_name_re,omitempty"`
	OnStep   *int    `json:"on_step,omitempty"    yaml:"on_step,omitempty"`
	WithFn   string  `json:"with_fn"              yaml:"with_fn"`
}

type Request struct {
	Compensations *[]requestCompensations `json:"compensations,omitempty"`
	Plan          string                  `json:"plan"`

	Init giraffe.Datum
}

type requestRead struct {
	Compensations *[]requestCompensations `json:"compensations,omitempty"`
	Plan          string                  `json:"plan"`
	Init          json.RawMessage         `json:"init"`
}

var _ serdes.Serde[Request] = (*requestSerde)(nil)

type requestSerde struct {
	datumSerde serdes.Serde[giraffe.Datum]
}

func (s requestSerde) Write(v Request) ([]byte, error) {
	//nolint:musttag
	raw, err := gson.Marshal(v)
	if err != nil {
		return nil, err
	}

	rawDat, err := s.datumSerde.Write(v.Init)
	if err != nil {
		return nil, E(err)
	}

	// Forgive me
	raw = raw[:1]
	raw = append(raw, []byte(`"init"`)...)
	raw = append(raw, rawDat...)
	raw = append(raw, '}')
	return raw, nil
}

func (s requestSerde) Read(b []byte) (Request, error) {
	read, err := gson.Unmarshal[requestRead](b)
	if err != nil {
		return Request{}, err
	}

	dat, err := s.datumSerde.Read(read.Init)
	if err != nil {
		return Request{}, E(err)
	}

	return Request{
		Compensations: read.Compensations,
		Plan:          read.Plan,
		Init:          dat,
	}, nil
}

func (s requestSerde) StreamTo(w io.Writer, v Request) error {
	enc := json.NewEncoder(w)
	//nolint:musttag
	return E(enc.Encode(v))
}

func (s requestSerde) StreamFrom(r io.Reader) (Request, error) {
	// TODO optimize this horrible impl.

	payload := new(bytes.Buffer)
	if _, err := io.Copy(payload, r); err != nil {
		return Request{}, E(err)
	}

	return s.Read(payload.Bytes())
}

func RequestSerde() serdes.Serde[Request] {
	return requestSerde{
		datumSerde: giraffe.DatumSerde(),
	}
}
