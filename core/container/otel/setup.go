package otel

import (
	"context"

	"github.com/hkoosha/giraffe/core/t11y/glog"
	"go.opentelemetry.io/otel/metric"

	"github.com/hkoosha/giraffe/core/container/finalizers"
	"github.com/hkoosha/giraffe/core/container/otel/internal"
	"github.com/hkoosha/giraffe/core/container/setup"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
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
)

func SetupOtel(namespace string) {
	setup.Finish("giraffe", "o11y", "setup")

	internal.Setup(namespace)

	invalidMetricOpCnt = M(
		internal.DefaultProvider().
			Meter("giraffe").
			Int64Counter("invalid_op"),
	)
}

func SetupOtelNoop() {
	internal.SetupNoop()
}

func Shutdown(ctx context.Context) {
	internal.Shutdown(ctx)
}

func Finalize(
	ctx context.Context,
) {
	fin.Finalize(ctx)
}
