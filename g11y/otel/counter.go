package otel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

type Int64Counter interface {
	Inc(context.Context, ...attribute.KeyValue)
}

type HTTPCounter interface {
	Inc(
		ctx context.Context,
		httpStatus int,
	)

	// IncErr treats http.StatusOk as error too (i.e., it as corrupt data in
	// the response's body).
	IncErr(
		ctx context.Context,
		httpStatus int,
	)

	NetworkFail(context.Context)

	UnexpectedResponseData(context.Context)
}

type OkCounter interface {
	Ok(context.Context, ...attribute.KeyValue)

	Fail(context.Context, ...attribute.KeyValue)
}

type HitOrMissCounter interface {
	Hit(context.Context, ...attribute.KeyValue)

	Miss(context.Context, ...attribute.KeyValue)
}
