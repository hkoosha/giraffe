package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	oprometheus "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/typing"
)

var shutdown func(context.Context)

func Setup(namespace string) {
	if !typing.IsSimpleMachineReadableName(namespace, 1, 32) {
		panic(EF("invalid namespace: %s", namespace))
	}

	exporter := M(oprometheus.New(
		oprometheus.WithNamespace(namespace),
		oprometheus.WithRegisterer(prometheus.DefaultRegisterer),
	))

	otel.SetMeterProvider(
		metric.NewMeterProvider(
			metric.WithReader(exporter),
		),
	)
}

func SetupNoop() {
	shutdown = func(context.Context) {}
}

func Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	shutdown(ctx)
}
