package remote

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/serdes"
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
	Init          giraffe.Datum           `json:"init"`
}

func RequestSerde() serdes.Serde[Request] {
	return serdes.Json[Request]()
}
