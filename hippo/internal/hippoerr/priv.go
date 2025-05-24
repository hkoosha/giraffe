package hippoerr

import (
	"reflect"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

type ErrorState interface {
	String(hE *HippoError) string
}

// NewHippoError Private function, do not call outside hippo package.
func NewHippoError(
	code int,
	msg string,
	state any,
) error {
	var errState ErrorState
	switch state.(type) {
	case *FnMissingKeysErrorState,
		*PipelineErrorState,
		*PlanErrorState,
		*RemoteErrorState:
		//nolint:errcheck,forcetypeassert
		errState = state.(ErrorState)

	default:
		panic(EF("unknown error type: %s", reflect.TypeOf(state).String()))
	}

	return E(&HippoError{
		code:       code,
		msg:        msg,
		errorState: errState,
	})
}
