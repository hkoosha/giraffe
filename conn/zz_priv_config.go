package conn

import (
	"net/http"
	"slices"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/hkoosha/giraffe/g11y"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

func (c *config) ensure() *config {
	g11y.NonNil(c)

	if !c.sealed {
		panic(EF("invalid config, did you use constructor to create one?"))
	}

	return c
}

func (c *config) open() *config {
	c.ensure()

	cp := *c
	return &cp
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
	pathPrefix string
	endpoint   string
	timeout    time.Duration
}

func (c *httpConfig) shallow() *httpConfig {
	return &httpConfig{
		pathPrefix: c.pathPrefix,
		endpoint:   c.endpoint,
		timeout:    c.timeout,
	}
}

func mkHttpConfig(
	timeout time.Duration,
) *httpConfig {
	if timeout < 1*time.Millisecond {
		panic("invalid timeout")
	}

	return &httpConfig{
		pathPrefix: "",
		endpoint:   "",
		timeout:    timeout,
	}
}

type respConfig struct {
	expectStatusCode   int
	expectNonEmptyBody bool
}

func (c *respConfig) shallow() *respConfig {
	return &respConfig{
		expectStatusCode:   c.expectStatusCode,
		expectNonEmptyBody: c.expectNonEmptyBody,
	}
}

func mkResponseConfig() *respConfig {
	return &respConfig{
		expectStatusCode:   -1,
		expectNonEmptyBody: false,
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
