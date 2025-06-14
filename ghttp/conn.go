package ghttp

import (
	"context"
	"strconv"

	"github.com/hkoosha/giraffe/zebra/serdes"
)

const (
	ReasonUnexpectedStatusCode FailureReason = 2
	ReasonEmptyResponse        FailureReason = 3
)

// =============================================================================

type FailureReason uint

type FailedResponseError struct {
	Reason FailureReason
	Resp   ConnResponse
}

func (e *FailedResponseError) Error() string {
	return "http request failed: " + strconv.Itoa(int(e.Reason))
}

// =============================================================================

type ConnResponse interface {
}

type Conn[T any, U any] interface {
	Call(
		ctx context.Context,
		method string,
		body T,
		path ...string,
	) (U, error)

	Patch(
		ctx context.Context,
		body T,
		path ...string,
	) (U, error)

	Put(
		ctx context.Context,
		body T,
		path ...string,
	) (U, error)

	Post(
		ctx context.Context,
		body T,
		path ...string,
	) (U, error)

	Get(
		ctx context.Context,
		path ...string,
	) (U, error)

	Delete(
		ctx context.Context,
		path ...string,
	) (U, error)
}

// =============================================================================

func NewConn[T, U any](
	cfg Config,
	tSerde serdes.Serde[T],
	uSerde serdes.Serde[U],
) Conn[T, U] {
	cloned := cfgOf(cfg)

	return newConn[T, U](cloned, tSerde, uSerde)
}
