package hippoerr

type remoteErrorState struct {
	err error
}

func (e *remoteErrorState) String(
	hE *hippoError,
) string {
	return "TODO::remoteErrorState :: " + hE.msg + " -> " + e.err.Error()
}

// NewRemoteError Private function, do not call outside hippo package.
func NewRemoteError(
	msg string,
	err error,
) error {
	return NewHippoError(
		ErrCodeRemoteCallFailure,
		msg,
		&remoteErrorState{
			err: err,
		},
	)
}
