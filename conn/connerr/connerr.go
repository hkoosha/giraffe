package connerr

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type UnhappyResponseError struct {
	ClosedResp *http.Response
	Reason     string
	Extra      string
	Code       int
}

func (e *UnhappyResponseError) Error() string {
	const sep = " => "
	const prefix = "unhappy response: "

	sb := new(strings.Builder)

	sb.Grow(len(prefix) + 3)
	sb.WriteString(prefix)
	sb.WriteString(strconv.Itoa(e.Code))

	if e.Reason != "" {
		sb.Grow(len(e.Reason) + len(sep))
		sb.WriteString(sep)
		sb.WriteString(e.Reason)
	}

	return sb.String()
}

func OfNil() error {
	return &UnhappyResponseError{
		Code:       -1,
		Reason:     "nil response",
		Extra:      "",
		ClosedResp: nil,
	}
}

func Of(resp *http.Response) error {
	if resp == nil {
		return OfNil()
	}

	unhappy := UnhappyResponseError{
		Code:       resp.StatusCode,
		Reason:     resp.Status,
		Extra:      "",
		ClosedResp: resp,
	}

	var err error
	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		if body, err0 := io.ReadAll(resp.Body); err0 != nil {
			unhappy.Extra = "!! could not read body => " + err0.Error()
			err = err0
		} else {
			unhappy.Extra = string(body)
		}
	}

	return errors.Join(&unhappy, err)
}

func OfEmpty(resp *http.Response) error {
	if resp == nil {
		return OfNil()
	}

	unhappy := UnhappyResponseError{
		Code:       resp.StatusCode,
		Reason:     resp.Status,
		Extra:      "!! empty response",
		ClosedResp: resp,
	}

	return &unhappy
}
