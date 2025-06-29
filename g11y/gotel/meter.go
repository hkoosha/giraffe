package gotel

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	"github.com/hkoosha/giraffe/g11y/finalizers"
	"github.com/hkoosha/giraffe/g11y/gotel/internal/metrics"
	"github.com/hkoosha/giraffe/g11y/setup"
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/typing"
)

func NewMetricBuilder(
	domain string,
	service string,
) *MetricBuilder {
	domainValid := typing.IsMachineReadableName(domain, 1, 32)
	svcValid := typing.IsMachineReadableName(service, 1, 32)

	if !domainValid || !svcValid {
		panic(EF("%s", "invalid domain or service name for metrics, domain="+
			domain+", service="+service))
	}

	name := domain + "/" + service

	return &MetricBuilder{
		domain:  domain,
		service: service,
		name:    name,
		isNoop:  false,
		meter:   nil,
		mu:      sync.Mutex{},
	}
}

type MetricBuilder struct {
	meter   metric.Meter
	domain  string
	service string
	name    string
	mu      sync.Mutex
	isNoop  bool
}

func (m *MetricBuilder) getMeter() metric.Meter {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.meter == nil {
		switch {
		case m.isNoop:
			m.meter = metrics.NoopProvider().Meter(m.name)

		default:
			m.meter = otel.GetMeterProvider().Meter(m.name)
		}
	}

	return m.meter
}

func (m *MetricBuilder) register(
	names []string,
	touch func(ctx context.Context),
) {
	what := "metric"
	if m.isNoop {
		what = "noop_" + what
	}

	for _, name := range names {
		setup.Finish("boot", "o11y", what, name)
	}

	finalizers.AddTo(fin, touch)
}

func (m *MetricBuilder) Counter(
	name string,
	description string,
) Int64Counter {
	cnt := metrics.NewCounter(
		m.getMeter(),
		m.domain,
		m.service,
		name,
		description,
		onInvalidMetric,
	)
	m.register(cnt.Once(), cnt.Touch)

	return cnt
}

func (m *MetricBuilder) OkCounter(
	name string,
	description string,
) OkCounter {
	cnt := metrics.NewOkCounter(
		m.getMeter(),
		m.domain,
		m.service,
		name,
		description,
		onInvalidMetric,
		"ok",
	)
	m.register(cnt.Once(), cnt.Touch)

	return cnt
}

func (m *MetricBuilder) HitOrMissCounter(
	name string,
	description string,
	label string,
) HitOrMissCounter {
	cnt := metrics.NewHitOrMissCounter(
		m.getMeter(),
		m.domain,
		m.service,
		name,
		description,
		onInvalidMetric,
		label,
	)
	m.register(cnt.Once(), cnt.Touch)

	return cnt
}

func (m *MetricBuilder) HTTPCounter(
	name string,
	description string,
) HTTPCounter {
	cnt := metrics.NewHTTPCounter(
		m.getMeter(),
		m.domain,
		m.service,
		name,
		description,
		onInvalidMetric,
	)
	m.register(cnt.Once(), cnt.Touch)

	return cnt
}

func (m *MetricBuilder) AsNoop() *MetricBuilder {
	if m.isNoop {
		return m
	}

	builder := NewMetricBuilder(m.domain, m.service)
	builder.isNoop = true

	return builder
}
