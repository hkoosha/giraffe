package hippoerr

import (
	"github.com/hkoosha/giraffe"
)

type fnErrorState struct {
	arg  giraffe.Datum
	keys []giraffe.Query
}

func (e *fnErrorState) String(
	hE *hippoError,
) string {
	_ = e.arg
	_ = e.keys
	return "TODO::fnErrorState :: " + hE.msg
}

// NewFnMissingKeysError Private function, do not call outside hippo package.
func NewFnMissingKeysError(
	arg giraffe.Datum,
	keys []giraffe.Query,
) error {
	return NewHippoError(
		ErrCodeMissingKeys,
		"missing keys",
		&fnErrorState{
			arg:  arg,
			keys: keys,
		},
	)
}
