package internal

import (
	"context"
	"net/http"
	"strings"

	"github.com/hkoosha/giraffe/core/t11y"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

func init() {
	for s := range 1023 {
		if http.StatusText(s) != "" {
			httpStatusAttributes[s] = attribute.Int("giraffe_http_status_code", s)
		}
	}
}

func NewCounter(
	meter metric.Meter,
	domain string,
	service string,
	name string,
	description string,
	onInvalidOp func(ctx context.Context, details string),
) *Counter {
	return &Counter{
		cnt: newBaseInt64Counter(meter, domain, service, name, description, onInvalidOp),
	}
}

type Counter struct {
	cnt baseInt64Counter
}

func (c *Counter) Touch(ctx context.Context) {
	c.cnt.touch(ctx)
}

func (c *Counter) Once() []string {
	return []string{c.cnt.fullName}
}

func (c *Counter) Inc(
	ctx context.Context,
	attrs ...attribute.KeyValue,
) {
	c.cnt.inc(ctx, attrs...)
}

// ==============================================================================.

var (
	httpStatusAttributeUnexpectedResponseData = attribute.Int("giraffe_http_status_code", -2)
	httpStatusAttributeNetworkFail            = attribute.Int("giraffe_http_status_code", -1)
	httpStatusAttributeInvalid                = attribute.Int("giraffe_http_status_code", 0)
	httpStatusAttributes                      = make(map[int]attribute.KeyValue)
)

func NewHTTPCounter(
	meter metric.Meter,
	domain string,
	service string,
	name string,
	description string,
	onInvalidOp func(ctx context.Context, details string),
) *HTTPCounter {
	mkCounter := func(extra ...attribute.KeyValue) baseInt64Counter {
		return newBaseInt64Counter(
			meter,
			domain,
			service,
			name,
			description,
			onInvalidOp,
			extra...,
		)
	}

	cnt := &HTTPCounter{
		toTouch: mkCounter(httpStatusAttributes[http.StatusOK]),
		cnt:     mkCounter(),
	}

	return cnt
}

type HTTPCounter struct {
	// Set to label status=ok.
	toTouch baseInt64Counter

	cnt baseInt64Counter
}

func (c *HTTPCounter) Once() []string {
	return []string{c.toTouch.fullName}
}

func (c *HTTPCounter) Touch(ctx context.Context) {
	c.toTouch.touch(ctx)
}

func (c *HTTPCounter) Inc(
	ctx context.Context,
	httpStatus int,
) {
	var a attribute.KeyValue
	if attr, ok := httpStatusAttributes[httpStatus]; ok {
		a = attr
	} else {
		a = httpStatusAttributeInvalid
	}

	c.cnt.inc(ctx, a)
}

func (c *HTTPCounter) NetworkFail(ctx context.Context) {
	c.cnt.inc(ctx, httpStatusAttributeNetworkFail)
}

func (c *HTTPCounter) UnexpectedResponseData(ctx context.Context) {
	c.cnt.inc(ctx, httpStatusAttributeUnexpectedResponseData)
}

func (c *HTTPCounter) IncErr(
	ctx context.Context,
	httpStatus int,
) {
	if 200 <= httpStatus && httpStatus <= 299 {
		c.UnexpectedResponseData(ctx)
	} else {
		c.Inc(ctx, httpStatus)
	}
}

// ==============================================================================.

func NewOkCounter(
	meter metric.Meter,
	domain string,
	service string,
	name string,
	description string,
	onInvalidOp func(ctx context.Context, details string),
	label string,
) *OkCounter {
	if !strings.HasPrefix(label, "giraffe") {
		label = "giraffe_" + label
	}

	if !t11y.IsMachineReadableName(label, 1, 32) {
		panic(EF("%s", "invalid label, "+
			"label="+label+", domain="+domain+", service="+service))
	}

	return &OkCounter{
		either: newEitherCounter(
			meter,
			domain,
			service,
			name,
			description,
			onInvalidOp,
			attribute.Bool(label, true),
			attribute.Bool(label, false),
		),
	}
}

type OkCounter struct {
	either *eitherCounter
}

func (c *OkCounter) Once() []string {
	return []string{c.either.left.fullName}
}

func (c *OkCounter) Touch(ctx context.Context) {
	c.either.left.touch(ctx)
	c.either.right.touch(ctx)
}

func (c *OkCounter) Ok(
	ctx context.Context,
	attr ...attribute.KeyValue,
) {
	c.either.incLeft(ctx, attr...)
}

func (c *OkCounter) Fail(
	ctx context.Context,
	attr ...attribute.KeyValue,
) {
	c.either.incRight(ctx, attr...)
}

// ==============================================================================.

func NewHitOrMissCounter(
	meter metric.Meter,
	domain string,
	service string,
	name string,
	description string,
	onInvalidOp func(ctx context.Context, details string),
	label string,
) *HitOrMissCounter {
	if !strings.HasPrefix(label, "giraffe") {
		label = "giraffe_" + label
	}

	if !t11y.IsMachineReadableName(label, 1, 32) {
		panic(EF("%s", "invalid label, "+
			"label="+label+", domain="+domain+", service="+service))
	}

	return &HitOrMissCounter{
		either: newEitherCounter(
			meter,
			domain,
			service,
			name,
			description,
			onInvalidOp,
			attribute.Bool("giraffe_hit", true),
			attribute.Bool("giraffe_hit", false),
		),
	}
}

type HitOrMissCounter struct {
	either *eitherCounter
}

func (c *HitOrMissCounter) Once() []string {
	return []string{c.either.left.fullName}
}

func (c *HitOrMissCounter) Touch(ctx context.Context) {
	c.either.left.touch(ctx)
	c.either.right.touch(ctx)
}

func (c *HitOrMissCounter) Hit(
	ctx context.Context,
	attr ...attribute.KeyValue,
) {
	c.either.incLeft(ctx, attr...)
}

func (c *HitOrMissCounter) Miss(
	ctx context.Context,
	attr ...attribute.KeyValue,
) {
	c.either.incRight(ctx, attr...)
}
