package connerrors

import (
	"errors"
)

var errEmptyResponse = errors.New("unexpected empty response")

func ErrEmptyResponse() error {
	return errEmptyResponse
}
