package gotel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/hkoosha/giraffe/g11y/finalizers"
	"github.com/hkoosha/giraffe/g11y/glog"
	"github.com/hkoosha/giraffe/g11y/gotel/internal/metrics"
	"github.com/hkoosha/giraffe/g11y/setup"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

var (
	fin = finalizers.NewFinalizer(setup.New())

	// Should be directly of underlying counter type, not our custom types to prevent recursion.
	invalidMetricOpCnt metric.Int64Counter

	onInvalidMetric = func(ctx context.Context, details string) {
		if invalidMetricOpCnt == nil {
			glog.Global().Error("invalid metric operation", details)
		} else {
			invalidMetricOpCnt.Add(ctx, 1)
		}
	}

	tracer trace.Tracer
)

func SetupOtel(
	name string,
	version string,
	ref string,
	namespace string,
) {
	setup.Finish("giraffe", "o11y", "setup")

	metrics.Setup(namespace)

	invalidMetricOpCnt = M(
		metrics.DefaultProvider().
			Meter("giraffe").
			Int64Counter("invalid_op"),
	)

	tracer = otel.Tracer(
		name,
		trace.WithInstrumentationAttributes(
			attribute.String("giraffe.version", version),
			attribute.String("giraffe.ref", ref),
		),
	)
}

func SetupOtelNoop() {
	metrics.SetupNoop()
}

func Shutdown(ctx context.Context) {
	metrics.Shutdown(ctx)
}

func Finalize(
	ctx context.Context,
) {
	fin.Finalize(ctx)
}
