package ghttp

import (
	"context"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/ghttp/headers"
	"github.com/hkoosha/giraffe/glog"
	"github.com/hkoosha/giraffe/zebra/z"
)

const (
	UserAgent = "Giraffe/1.0"
)

type RetryIfFn = func(
	ctx context.Context,
	resp *http.Response,
	err error,
	attempt uint,
	cfg *Config,
) (bool, error)

func NewConfig(
	lg glog.GLog,
	timeout time.Duration,
) *Config {
	cfg := &Config{
		base:   defaultTransport,
		lg:     lg,
		rt:     nil,
		seal_:  seal{false},
		resp:   mkResponseConfig(),
		http:   mkHttpConfig(timeout),
		header: mkHeaderConfig(),
		log:    mkLogConfig(),
		retry:  mkRetryConfig(),
		otel:   mkOtelConfig(),
	}

	cfg.seal()

	return cfg
}

// Config
// Keep all setter methods prefixed with either of:
// - Is
// - With
// - Without
// - Set
// - And
// This makes refactoring easier and more consistent, since all the methods
// minus the getters are implemented in [Client] too.
type Config struct {
	lg     glog.GLog //nolint:unused
	base   http.RoundTripper
	rt     http.RoundTripper
	resp   *respConfig
	http   *httpConfig
	header *headerConfig
	log    *logConfig
	retry  *retryConfig
	otel   *otelConfig
	seal_  seal
}

// =============================================================================.

func (c *Config) Ensure() {
	c.ensure()
}

func (c *Config) Std() *http.Client {
	return &http.Client{
		Transport:     c.mkTransport(),
		Timeout:       c.http.timeout,
		CheckRedirect: nil,
		Jar:           nil,
	}
}

func (c *Config) WithLg(lg glog.GLog) *Config {
	g11y.NonNil(lg)

	cp := c.open()
	cp.lg = lg
	cp.seal()

	return cp
}

func (c *Config) IsLogged() bool {
	return c.ensure().log.isLogged
}

func (c *Config) WithLogged() *Config {
	return c.SetLogged(true)
}

func (c *Config) WithoutLogged() *Config {
	return c.SetLogged(false)
}

func (c *Config) SetLogged(b bool) *Config {
	if c.log.isLogged == b {
		return c
	}

	cp := c.open()
	cp.log = cp.log.shallow()
	cp.log.isLogged = b
	cp.seal()

	return cp
}

func (c *Config) IsPlainLog() bool {
	return c.ensure().log.isPlainLog
}

func (c *Config) WithPlainLog() *Config {
	return c.SetPlainLog(true)
}

func (c *Config) WithoutPlainLog() *Config {
	return c.SetPlainLog(false)
}

func (c *Config) SetPlainLog(b bool) *Config {
	if c.log.isPlainLog == b {
		return c
	}

	cp := c.open()
	cp.log = cp.log.shallow()
	cp.log.isPlainLog = b
	cp.seal()

	return cp
}

func (c *Config) LgHeaderFilter() HeaderFilter {
	return c.ensure().log.headerFilter
}

func (c *Config) WithLgFilteredHeaders(f HeaderFilter) *Config {
	g11y.NonNil(f)

	cp := c.open()
	cp.log = cp.log.shallow()
	cp.log.headerFilter = f
	cp.seal()

	return cp
}

func (c *Config) LgHeaderMask() HeaderFilter {
	return c.ensure().log.maskedHeaders
}

func (c *Config) WithLgMaskedHeaders(f HeaderFilter) *Config {
	g11y.NonNil(f)

	cp := c.open()
	cp.log = cp.log.shallow()
	cp.log.maskedHeaders = f
	cp.seal()

	return cp
}

func (c *Config) WithHeaderOverwrites(
	withDefaults bool,
	h map[string]string,
) *Config {
	g11y.NonNil(h)

	if withDefaults {
		h = z.UnionLeft(h, defaultHeaders)
	} else {
		h = maps.Clone(h)
	}

	if z.MapEq(c.header.overwrite, h) {
		return c
	}

	cp := c.open()
	cp.header = cp.header.shallow()
	cp.header.overwrite = h
	cp.seal()

	return cp
}

func (c *Config) WithBearerToken(bt string) *Config {
	if strings.TrimSpace(bt) == "" {
		panic("empty bearer token")
	}

	bt = withBearerPrefix(bt)

	if c.header.overwrite[headers.Authorization] == bt {
		return c
	}

	cp := c.open()
	cp.header = cp.header.shallow()
	cp.header.overwrite = maps.Clone(cp.header.overwrite)
	cp.header.overwrite[headers.Authorization] = bt
	cp.seal()
	return cp
}

func (c *Config) WithBearerProvider(fn HeaderFn) *Config {
	g11y.NonNil(fn)

	fn = func(ctx context.Context, config *Config) string {
		return withBearerPrefix(fn(ctx, config))
	}

	cp := c.open()
	cp.header = cp.header.shallow()
	cp.header.overwriteFns = maps.Clone(cp.header.overwriteFns)
	cp.header.overwriteFns[headers.Authorization] = fn
	cp.seal()

	return cp
}

func (c *Config) WithoutBearerProvider() *Config {
	_, ok := c.header.overwriteFns[headers.Authorization]
	if !ok {
		return c
	}

	cp := c.open()
	cp.header = cp.header.shallow()
	cp.header.overwriteFns = maps.Clone(cp.header.overwriteFns)
	delete(cp.header.overwriteFns, headers.Authorization)
	cp.seal()

	return cp
}

func (c *Config) IsExpecting2xx() bool {
	return c.ensure().resp.isExpecting2xx
}

func (c *Config) WithExpecting2xx() *Config {
	return c.SetExpecting2xx(true)
}

func (c *Config) WithoutExpecting2xx() *Config {
	return c.SetExpecting2xx(false)
}

func (c *Config) SetExpecting2xx(b bool) *Config {
	if c.resp.isExpecting2xx == b {
		return c
	}

	cp := c.open()
	cp.resp = cp.resp.shallow()
	cp.resp.isExpecting2xx = b
	cp.seal()

	return cp
}

func (c *Config) Endpoint() string {
	return c.ensure().http.endpoint
}

func (c *Config) WithEndpoint(e string) *Config {
	if e == "" || strings.TrimSpace(e) == "" {
		panic("empty endpoint")
	}

	if c.http.endpoint == e {
		return c
	}

	cp := c.open()
	cp.http = cp.http.shallow()
	cp.http.endpoint = e
	cp.seal()

	return cp
}

func (c *Config) WithoutEndpoint() *Config {
	if c.http.endpoint == "" {
		return c
	}

	cp := c.open()
	cp.http = cp.http.shallow()
	cp.http.endpoint = ""
	cp.seal()

	return cp
}

func (c *Config) PathPrefix() string {
	return c.ensure().http.pathPrefix
}

func (c *Config) WithPathPrefix(p string) *Config {
	// TODO validate p.
	if p == "" || strings.TrimSpace(p) == "" {
		panic("empty path prefix")
	}

	if c.http.pathPrefix == p {
		return c
	}

	cp := c.open()
	cp.http = cp.http.shallow()
	cp.http.pathPrefix = p
	cp.seal()

	return cp
}

func (c *Config) WithoutPathPrefix() *Config {
	if c.http.pathPrefix == "" {
		return c
	}

	cp := c.open()
	cp.http = cp.http.shallow()
	cp.http.pathPrefix = ""
	cp.seal()

	return cp
}

func (c *Config) AndPathPrefix(p string) *Config {
	if p == "" || strings.TrimSpace(p) == "" {
		panic("empty path prefix")
	}

	cp := c.WithPathPrefix(Join(c.http.pathPrefix, p))
	cp.seal()

	return cp
}

func (c *Config) Timeout() time.Duration {
	return c.ensure().http.timeout
}

func (c *Config) WithTimeout(t time.Duration) *Config {
	if c.http.timeout == t {
		return c
	}

	cp := c.open()
	cp.http = cp.http.shallow()
	cp.http.timeout = t
	cp.seal()

	return cp
}

func (c *Config) WithTransport(t http.RoundTripper) *Config {
	g11y.NonNil(t)

	cp := c.open()
	cp.base = t
	cp.seal()

	return cp
}

func (c *Config) WithTraced() *Config {
	return c.SetTraced(true)
}

func (c *Config) WithoutTraced() *Config {
	return c.SetTraced(false)
}

func (c *Config) SetTraced(b bool) *Config {
	if c.otel.enabled == b {
		return c
	}

	cp := c.open()
	cp.otel = cp.otel.shallow()
	cp.otel.enabled = b
	cp.seal()

	return cp
}

func (c *Config) WithTraceOptions(
	options ...otelhttp.Option,
) *Config {
	if slices.Equal(c.otel.options, options) {
		return c
	}

	cp := c.open()
	cp.otel = cp.otel.shallow()
	cp.otel.options = options
	cp.seal()

	return cp
}

func (c *Config) IsTraced() bool {
	return c.otel.enabled
}

// =============================================================================.

func (c *Config) IsLogReties() bool {
	return c.ensure().retry.logged
}

func (c *Config) WithLogReties() *Config {
	return c.SetLogReties(true)
}

func (c *Config) WithoutLogReties() *Config {
	return c.SetLogReties(false)
}

func (c *Config) SetLogReties(b bool) *Config {
	if c.retry.logged == b {
		return c
	}

	cp := c.open()
	cp.retry = cp.retry.shallow()
	cp.retry.logged = b
	cp.seal()

	return cp
}

func (c *Config) RetryMax() uint {
	return c.ensure().retry.maxRetries
}

func (c *Config) WithMaxRetries(r uint) *Config {
	if c.retry.maxRetries == r {
		return c
	}

	cp := c.open()
	cp.retry = cp.retry.shallow()
	cp.retry.maxRetries = r
	cp.seal()

	return cp
}

func (c *Config) RetriedStatusCodes() []int {
	return slices.Clone(c.ensure().retry.retryIfStatuses)
}

func (c *Config) WithRetriedStatusCodes(sc ...int) *Config {
	g11y.NonNil(sc)

	if slices.Equal(c.retry.retryIfStatuses, sc) {
		return c
	}

	cp := c.open()
	cp.retry = cp.retry.shallow()
	cp.retry.retryIfStatuses = slices.Clone(sc)
	cp.seal()

	return cp
}
