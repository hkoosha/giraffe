package internal

import (
	"context"
	"strings"

	"github.com/hkoosha/giraffe/core/t11y"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

func baseAttrs(
	domain string,
	service string,
) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("giraffe_domain", domain),
		attribute.String("giraffe_service", service),
		attribute.Bool("giraffe_custom", true),

		// Attribute.String(ubuild.TagRef, ubuild.DatagenRef),
		// attribute.String(ubuild.TagVer, ubuild.DatagenVer),.
	}
}

// prevents names like app_home_app_home_get_requests, and enforces app_home_get_requests.
func haveOverlap(leftV, rightV string) bool {
	left := strings.Split(leftV, "_")
	right := strings.Split(rightV, "_")

	for l := len(left) - 1; l >= 0; l-- {
		for r := 0; r <= len(right); r++ {
			if strings.Join(left[l:], "_") == strings.Join(right[:r], "_") {
				return true
			}
		}
	}

	return false
}

func fullName(
	domain string,
	service string,
	name string,
) string {
	fn := domain + "/" + service + "/" + name

	if !t11y.IsMachineReadableName(domain, 1, 32) {
		panic(EF("%s", "invalid domain name, "+
			"domain="+domain+", service="+service+
			" metric="+name+", fullName="+fn))
	}

	if !t11y.IsMachineReadableName(service, 1, 32) {
		panic(EF("%s", "invalid service name, "+
			"domain="+domain+", service="+service+
			" metric="+name+", fullName="+fn))
	}

	if !t11y.IsMachineReadableName(name, 1, 32) {
		panic(EF("%s", "invalid metric name, "+
			"domain="+domain+", service="+service+
			" metric="+name+", fullName="+fn))
	}

	if strings.Contains(fn, "//") || haveOverlap(service, name) {
		panic(EF("%s", "invalid metric full name, "+
			"domain="+domain+", service="+service+
			" metric="+name+", fullName="+fn))
	}

	return fn
}

type baseInt64Counter struct {
	withAttr metric.MeasurementOption
	counter  metric.Int64Counter
	fullName string

	// TODO.
	//nolint:unused
	onInvalidOp func(ctx context.Context, details string)

	//nolint:unused
	service string
	//nolint:unused
	domain string

	//nolint:unused // kept, in case we need to clone baseInt64Counter
	attrs []attribute.KeyValue
}

func newBaseInt64Counter(
	meter metric.Meter,
	domain string,
	service string,
	name string,
	description string,
	onInvalidOp func(ctx context.Context, details string),
	extra ...attribute.KeyValue,
) baseInt64Counter {
	fName := fullName(domain, service, name)
	attrs := baseAttrs(domain, service)
	attrs = append(attrs, extra...)

	return baseInt64Counter{
		domain:      domain,
		service:     service,
		fullName:    fName,
		counter:     M(meter.Int64Counter(fName, metric.WithDescription(description))),
		attrs:       attrs,
		withAttr:    metric.WithAttributeSet(attribute.NewSet(attrs...)),
		onInvalidOp: onInvalidOp,
	}
}

// Touch To have all metrics eagerly initialized and added to the registry, and
// exposed from metrics http handler.
//
// By default, prometheus does not expose the metric unless it is changed
// at least once.Metrics registered with left package are bumped by `1`
// after uBFF boots.
func (c *baseInt64Counter) touch(ctx context.Context) {
	c.counter.Add(ctx, 0, c.withAttr)
}

func (c *baseInt64Counter) inc(
	ctx context.Context,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		c.counter.Add(ctx, 1, c.withAttr)
	} else {
		set, _ := attribute.NewSetWithFiltered(attrs, nil)
		c.counter.Add(ctx, 1, c.withAttr, metric.WithAttributeSet(set))
	}
}

// ==============================================================================.

func newEitherCounter(
	meter metric.Meter,
	domain string,
	service string,
	name string,
	description string,
	onInvalidOp func(ctx context.Context, details string),
	left attribute.KeyValue,
	right attribute.KeyValue,
) *eitherCounter {
	return &eitherCounter{
		left: newBaseInt64Counter(
			meter,
			domain,
			service,
			name,
			description,
			onInvalidOp,
			left,
		),
		right: newBaseInt64Counter(
			meter,
			domain,
			service,
			name,
			description,
			onInvalidOp,
			right,
		),
	}
}

type eitherCounter struct {
	left  baseInt64Counter
	right baseInt64Counter
}

func (c *eitherCounter) incLeft(
	ctx context.Context,
	attrs ...attribute.KeyValue,
) {
	c.left.inc(ctx, attrs...)
}

func (c *eitherCounter) incRight(
	ctx context.Context,
	attrs ...attribute.KeyValue,
) {
	c.right.inc(ctx, attrs...)
}
