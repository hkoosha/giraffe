package hippo

type ExternalError struct {
	err error
}

func (e *ExternalError) Error() string {
	return "EXTERNAL_ERROR: " + e.err.Error()
}
