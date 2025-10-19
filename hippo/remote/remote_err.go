package remote

import (
	"net/http"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

var (
	errUnknownError = &RemoteError{
		httpStatus:    http.StatusInternalServerError,
		userSafeError: ErrorPayload{"unknown error"},
	}

	errParsingPayload = &RemoteError{
		httpStatus:    http.StatusBadRequest,
		userSafeError: ErrorPayload{"error while parsing payload"},
	}

	errProcessingRequest = &RemoteError{
		httpStatus:    http.StatusUnprocessableEntity,
		userSafeError: ErrorPayload{"error while processing request"},
	}
)

type ErrorPayload struct {
	Message string `json:"message"`
}

//goland:noinspection GoNameStartsWithPackageName
type RemoteError struct {
	userSafeError any
	httpStatus    int
}

func (e *RemoteError) Error() string {
	return "remote error"
}

func (e *RemoteError) StatusCode() int {
	return e.httpStatus
}

func (e *RemoteError) UserSafeError() any {
	if e.userSafeError == nil {
		return errUnknownError.userSafeError
	}

	return e.userSafeError
}

// =============================================================================.

func newErrorParsingPayload(err error) error {
	return E(err, errParsingPayload)
}

func newErrorProcessingRequest(err error) error {
	return E(err, errProcessingRequest)
}

func newErrorMissingPlan(
	plan string,
) error {
	err := &RemoteError{
		httpStatus:    http.StatusBadRequest,
		userSafeError: ErrorPayload{"missing plan: " + plan},
	}

	return E(err)
}

func newErrorMissingFn(
	fn string,
) error {
	err := &RemoteError{
		httpStatus:    http.StatusBadRequest,
		userSafeError: ErrorPayload{"missing fn: " + fn},
	}

	return E(err)
}

func newUnknownError(err error) error {
	return E(err, errUnknownError)
}
