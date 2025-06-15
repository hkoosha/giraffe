package ghttp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/zebra/serdes"
)

const (
	ReasonUnexpectedStatusCode FailureReason = 2
	ReasonEmptyResponse        FailureReason = 3
)

// =============================================================================.

type FailureReason uint

type FailedResponseError struct {
	Resp   ConnResponse
	Reason FailureReason
}

func (e *FailedResponseError) Error() string {
	return "http request failed: " + strconv.FormatUint(uint64(e.Reason), 10)
}

// =============================================================================.

type ConnResponse any

type Conn[Q any, R any] interface {
	Std() *http.Client
	Cfg() Config

	Call(
		ctx context.Context,
		method string,
		body Q,
		path ...string,
	) (R, error)

	Patch(
		ctx context.Context,
		body Q,
		path ...string,
	) (R, error)

	Put(
		ctx context.Context,
		body Q,
		path ...string,
	) (R, error)

	Post(
		ctx context.Context,
		body Q,
		path ...string,
	) (R, error)

	Get(
		ctx context.Context,
		path ...string,
	) (R, error)

	Delete(
		ctx context.Context,
		path ...string,
	) (R, error)
}

// =============================================================================.

func NewConn[Q, R any](
	cfg Config,
	tSerde serdes.Serde[Q],
	uSerde serdes.Serde[R],
) Conn[Q, R] {
	cloned := cfgOf(cfg)

	return newConn[Q, R](cloned, tSerde, uSerde)
}

func NewJsonConn[Q, R any](
	cfg Config,
) Conn[Q, R] {
	return NewConn[Q, R](cfg, serdes.JsonSerde[Q](), serdes.JsonSerde[R]())
}

// =============================================================================.

func ToJsonType[T, U, Q, R any](
	conn Conn[T, U],
) Conn[Q, R] {
	g11y.NonNil(conn)

	cfg := conn.Cfg()
	return NewConn[Q, R](
		cfg,
		serdes.JsonSerde[Q](),
		serdes.JsonSerde[R](),
	)
}
