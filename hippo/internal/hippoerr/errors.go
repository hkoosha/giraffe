package hippoerr

const (
	ErrCodeFailedStep = iota + 1
	ErrCodeMissingKeys
	ErrCodeMissingFn
	ErrCodeDuplicateFn
	ErrCodeInvalidStepName
	ErrCodeRemoteCallFailure
)

type HippoError struct {
	errorState ErrorState
	msg        string
	code       int
}

func (e *HippoError) State() ErrorState {
	return e.errorState
}

func (e *HippoError) Error() string {
	return e.errorState.String(e)
}

func (e *HippoError) Code() int {
	return e.code
}
