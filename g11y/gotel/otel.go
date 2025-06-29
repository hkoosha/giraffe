package gotel

import (
	"go.opentelemetry.io/otel/trace"
)

func Tracer() trace.Tracer {
	return tracer
}
