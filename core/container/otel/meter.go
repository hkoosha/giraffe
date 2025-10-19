package otel

import (
	"context"
	"sync"

	"github.com/hkoosha/giraffe/core/t11y"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	"github.com/hkoosha/giraffe/core/container/finalizers"
	"github.com/hkoosha/giraffe/core/container/otel/internal"
	"github.com/hkoosha/giraffe/core/container/setup"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

func NewMetricBuilder(
	domain string,
	service string,
) *MetricBuilder {
	domainValid := t11y.IsMachineReadableName(domain, 1, 32)
	svcValid := t11y.IsMachineReadableName(service, 1, 32)

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
			m.meter = internal.NoopProvider().Meter(m.name)

		default:
			m.meter = otel.GetMeterProvider().Meter(m.name)
		}
	}

	return m.meter
}

func (m *MetricBuilder) register(names []string, touch func(ctx context.Context)) {
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
	cnt := internal.NewCounter(
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
	cnt := internal.NewOkCounter(
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
	cnt := internal.NewHitOrMissCounter(
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
	cnt := internal.NewHTTPCounter(
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
