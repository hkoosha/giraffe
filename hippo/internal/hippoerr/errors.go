package hippoerr

const (
	ErrCodeFailedStep = iota + 1
	ErrCodeMissingKeys
	ErrCodeMissingFn
	ErrCodeDuplicateFn
	ErrCodeInvalidStepName
	ErrCodeRemoteCallFailure
)

type hippoError struct {
	errorState errorState
	msg        string
	code       int
}

func (e *hippoError) Error() string {
	return e.errorState.String(e)
}

func (e *hippoError) Code() int {
	return e.code
}
