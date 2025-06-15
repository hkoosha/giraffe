package otel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

type Int64Counter interface {
	Inc(
		ctx context.Context,
		attrs ...attribute.KeyValue,
	)
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

	NetworkFail(ctx context.Context)

	UnexpectedResponseData(ctx context.Context)
}

type OkCounter interface {
	Ok(
		ctx context.Context,
		attr ...attribute.KeyValue,
	)

	Fail(
		ctx context.Context,
		attr ...attribute.KeyValue,
	)
}

type HitOrMissCounter interface {
	Hit(
		ctx context.Context,
		attr ...attribute.KeyValue,
	)

	Miss(
		ctx context.Context,
		attr ...attribute.KeyValue,
	)
}
