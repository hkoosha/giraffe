package hippoerr

import (
	"reflect"
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
		*FnMissingOutputErrorState,
		*PipelineErrorState,
		*PlanErrorState,
		*RemoteErrorState:
		//nolint:errcheck,forcetypeassert
		errState = state.(ErrorState)

	default:
		panic("unknown error type: " + reflect.TypeOf(state).String())
	}

	return &HippoError{
		code:       code,
		msg:        msg,
		errorState: errState,
	}
}
