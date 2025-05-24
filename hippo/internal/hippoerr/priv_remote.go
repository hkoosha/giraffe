package hippoerr

type RemoteErrorState struct {
	err error
}

func (e *RemoteErrorState) String(
	hE *HippoError,
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
		&RemoteErrorState{
			err: err,
		},
	)
}
