package hippoerr

import (
	"reflect"
)

type errorState interface {
	String(hE *hippoError) string
}

// NewHippoError Private function, do not call outside hippo package.
func NewHippoError(
	code int,
	msg string,
	state any,
) error {
	var errState errorState
	switch state.(type) {
	case *fnErrorState,
		*pipelineErrorState,
		*planErrorState,
		*remoteErrorState:
		//nolint:errcheck,forcetypeassert
		errState = state.(errorState)

	default:
		panic("unknown error type: " + reflect.TypeOf(state).String())
	}

	return &hippoError{
		code:       code,
		msg:        msg,
		errorState: errState,
	}
}
