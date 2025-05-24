package internal

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	mu                                = &sync.Mutex{}
	noopProvider metric.MeterProvider = nil
)

func NoopProvider() metric.MeterProvider {
	mu.Lock()
	defer mu.Unlock()

	if noopProvider == nil {
		noopProvider = noop.NewMeterProvider()
	}

	return noopProvider
}

func DefaultProvider() metric.MeterProvider {
	return otel.GetMeterProvider()
}
