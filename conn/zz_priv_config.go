package conn

import (
	"maps"
	"net/http"
	"slices"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/hkoosha/giraffe/core/t11y"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

func (c *config) ensure() *config {
	t11y.NonNil(c)

	if !c.sealed {
		panic(EF("invalid config, did you use constructor to create one?"))
	}

	return c
}

func (c *config) open() *config {
	c.ensure()

	return &*c
}

func (c *config) seal() {
	if c.sealed {
		return
	}

	var rt http.RoundTripper = &giraffeRT{
		cfg: c,
		rt:  c.base,
	}

	if c.otel.enabled {
		rt = otelhttp.NewTransport(
			rt,
			c.otel.options...,
		)
	}

	c.rt = rt

	c.sealed = true
}

func (c *config) mkTransport() http.RoundTripper {
	return giraffeRT{
		cfg: c,
		rt:  c.rt,
	}
}

// =============================================================================.

type retryConfig struct {
	retryIf         RetryIfFn
	retryIfStatuses []int
	maxRetries      uint
	backoffDuration time.Duration
	logged          bool
}

func (c *retryConfig) shallow() *retryConfig {
	return &retryConfig{
		logged:          c.logged,
		maxRetries:      c.maxRetries,
		backoffDuration: c.backoffDuration,
		retryIfStatuses: slices.Clone(c.retryIfStatuses),
		retryIf:         c.retryIf,
	}
}

func mkRetryConfig() *retryConfig {
	return &retryConfig{
		logged:          false,
		maxRetries:      3,
		backoffDuration: 1 * time.Second,
		retryIfStatuses: []int{
			http.StatusRequestTimeout,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},
		retryIf: nil,
	}
}

type logConfig struct {
	headerFilter  HeaderFilter
	maskedHeaders HeaderFilter
	isPlainLog    bool
	isLogged      bool
}

func (c *logConfig) shallow() *logConfig {
	return &logConfig{
		isPlainLog:    c.isPlainLog,
		isLogged:      c.isLogged,
		headerFilter:  c.headerFilter,
		maskedHeaders: c.maskedHeaders,
	}
}

func mkLogConfig() *logConfig {
	return &logConfig{
		isPlainLog:    false,
		isLogged:      false,
		headerFilter:  defaultHeaderFilter,
		maskedHeaders: defaultHeaderMasked,
	}
}

type otelConfig struct {
	rt      http.RoundTripper
	options []otelhttp.Option
	enabled bool
}

func (c *otelConfig) shallow() *otelConfig {
	return &otelConfig{
		rt:      c.rt,
		options: c.options,
		enabled: c.enabled,
	}
}

func mkOtelConfig() *otelConfig {
	return &otelConfig{
		rt:      nil,
		options: nil,
		enabled: false,
	}
}

type httpConfig struct {
	endpointsByName map[string]string
	defaultMethod   string
	pathPrefix      string
	endpoint        string
	timeout         time.Duration
}

func (c *httpConfig) shallow() *httpConfig {
	return &httpConfig{
		defaultMethod:   c.defaultMethod,
		pathPrefix:      c.pathPrefix,
		endpoint:        c.endpoint,
		timeout:         c.timeout,
		endpointsByName: maps.Clone(c.endpointsByName),
	}
}

func mkHttpConfig(
	timeout time.Duration,
) *httpConfig {
	if timeout < 1*time.Millisecond {
		panic(EF("invalid timeout"))
	}

	return &httpConfig{
		defaultMethod:   http.MethodGet,
		pathPrefix:      "",
		endpoint:        "",
		timeout:         timeout,
		endpointsByName: make(map[string]string),
	}
}

type respConfig struct {
	expectStatusCode []int
}

func (c *respConfig) clone() *respConfig {
	return &respConfig{
		expectStatusCode: slices.Clone(c.expectStatusCode),
	}
}

func mkResponseConfig() *respConfig {
	return &respConfig{
		expectStatusCode: nil,
	}
}

type headerConfig struct {
	overwrite   map[string]string
	overwriters map[string]HeaderProvider
}

func mkHeaderConfig() *headerConfig {
	return &headerConfig{
		overwrite:   nil,
		overwriters: nil,
	}
}

func (c *headerConfig) shallow() *headerConfig {
	return &headerConfig{
		overwrite:   c.overwrite,
		overwriters: c.overwriters,
	}
}
