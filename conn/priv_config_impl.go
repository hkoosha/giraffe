package conn

import (
	"context"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/hkoosha/giraffe/conn/headers"
	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/glog"
	"github.com/hkoosha/giraffe/zebra/serdes"
	"github.com/hkoosha/giraffe/zebra/z"
)

func cfgOf(
	cfg Config,
) *config {
	cfg.Ensure()

	if cast, ok := cfg.(*config); ok {
		return cast
	}

	cast := newConfig(cfg.Lg(), cfg.Timeout(), cfg.Serde()).
		setPlainLog(cfg.IsPlainLog()).
		setLogged(cfg.IsLogged()).
		setLogReties(cfg.IsRetryLog()).
		withMaxRetries(cfg.RetryMax()).
		setTraced(cfg.IsTraced()).
		withRetryBackoffDuration(cfg.RetryBackoffDuration())

	if v := cfg.ExpectingStatusCode(); v > 0 {
		cast = cast.withExpectingStatusCode(v)
	}
	if v := cfg.PathPrefix(); v != "" {
		cast = cast.withPathPrefix(v)
	}
	if v := cfg.Endpoint(); v != "" {
		cast = cast.withEndpoint(v)
	}
	if v := cast.LgHeaderFilter(); v != nil {
		cast = cast.withLgFilteredHeaders(v)
	}
	if v := cast.LgHeaderMask(); v != nil {
		cast = cast.withLgMaskedHeaders(v)
	}
	if v := cfg.HeaderOverwrites(); v != nil {
		cast = cast.withHeaderOverwrites(false, v)
	}
	if v := cfg.RetryStatusCodes(); v != nil {
		cast = cast.withRetriedStatusCodes(v...)
	}
	if v := cfg.TraceOptions(); v != nil {
		cast = cast.withTraceOptions(v...)
	}
	if v := cfg.RetryIf(); v != nil {
		cast = cast.withRetryIf(v)
	}

	for header := range cfg.HeaderOverwriters() {
		if header != headers.Authorization {
			panic("invalid config, unsupported header overwriter: " + header)
		}
	}
	cast = cast.withHeaderOverwriters(cfg.HeaderOverwriters())

	return cast
}

func newConfig(
	lg glog.Lg,
	timeout time.Duration,
	serde serdes.Serde[any],
) *config {
	cfg := &config{
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
		serde_: serde,
	}

	cfg.seal()

	return cfg
}

type config struct {
	lg     glog.Lg
	base   http.RoundTripper
	rt     http.RoundTripper
	serde_ serdes.Serde[any]
	resp   *respConfig
	http   *httpConfig
	header *headerConfig
	log    *logConfig
	retry  *retryConfig
	otel   *otelConfig
	seal_  seal
}

// =============================================================================.

func (c *config) Serde() serdes.Serde[any] {
	return c.serde()
}

func (c *config) serde() serdes.Serde[any] {
	c.ensure()

	return c.serde_
}

func (c *config) Conn() Conn[[]byte] {
	return c.connection()
}

func (c *config) connection() Conn[[]byte] {
	c.ensure()

	return newConn(c, serdes.Bytes())
}

func (c *config) Ensure() {
	c.ensure()
}

func (c *config) Lg() glog.Lg {
	return c.lg
}

func (c *config) Std() *http.Client {
	return &http.Client{
		Transport:     c.mkTransport(),
		Timeout:       c.http.timeout,
		CheckRedirect: nil,
		Jar:           nil,
	}
}

func (c *config) WithLg(lg glog.Lg) Config {
	return c.withLg(lg)
}

func (c *config) withLg(lg glog.Lg) *config {
	g11y.NonNil(lg)

	cp := c.open()
	cp.lg = lg
	cp.seal()

	return cp
}

func (c *config) IsLogged() bool {
	return c.ensure().log.isLogged
}

func (c *config) WithLogged() Config {
	return c.SetLogged(true)
}

func (c *config) WithoutLogged() Config {
	return c.withoutLogged()
}

func (c *config) withoutLogged() *config {
	return c.setLogged(false)
}

func (c *config) SetLogged(b bool) Config {
	return c.setLogged(b)
}

func (c *config) setLogged(b bool) *config {
	if c.log.isLogged == b {
		return c
	}

	cp := c.open()
	cp.log = cp.log.shallow()
	cp.log.isLogged = b
	cp.seal()

	return cp
}

func (c *config) IsPlainLog() bool {
	return c.ensure().log.isPlainLog
}

func (c *config) WithPlainLog() Config {
	return c.withPlainLog()
}

func (c *config) withPlainLog() *config {
	return c.setPlainLog(true)
}

func (c *config) WithoutPlainLog() Config {
	return c.withoutPlainLog()
}

func (c *config) withoutPlainLog() *config {
	return c.setPlainLog(false)
}

func (c *config) SetPlainLog(b bool) Config {
	return c.setPlainLog(b)
}

func (c *config) setPlainLog(b bool) *config {
	if c.log.isPlainLog == b {
		return c
	}

	cp := c.open()
	cp.log = cp.log.shallow()
	cp.log.isPlainLog = b
	cp.seal()

	return cp
}

func (c *config) LgHeaderFilter() HeaderFilter {
	return c.ensure().log.headerFilter
}

func (c *config) WithLgFilteredHeaders(
	f HeaderFilter,
) Config {
	return c.withLgFilteredHeaders(f)
}

func (c *config) withLgFilteredHeaders(
	f HeaderFilter,
) *config {
	g11y.NonNil(f)

	cp := c.open()
	cp.log = cp.log.shallow()
	cp.log.headerFilter = f
	cp.seal()

	return cp
}

func (c *config) LgHeaderMask() HeaderFilter {
	return c.ensure().log.maskedHeaders
}

func (c *config) WithLgMaskedHeaders(
	f HeaderFilter,
) Config {
	return c.withLgMaskedHeaders(f)
}

func (c *config) withLgMaskedHeaders(
	f HeaderFilter,
) *config {
	g11y.NonNil(f)

	cp := c.open()
	cp.log = cp.log.shallow()
	cp.log.maskedHeaders = f
	cp.seal()

	return cp
}

func (c *config) HeaderOverwrites() map[string]string {
	return maps.Clone(c.ensure().header.overwrite)
}

func (c *config) HeaderOverwriters() map[string]HeaderProvider {
	return maps.Clone(c.ensure().header.overwriters)
}

func (c *config) WithHeaderOverwrites(
	withDefaults bool,
	h map[string]string,
) Config {
	return c.withHeaderOverwrites(withDefaults, h)
}

func (c *config) withHeaderOverwrites(
	withDefaults bool,
	h map[string]string,
) *config {
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

func (c *config) withHeaderOverwriters(
	h map[string]HeaderProvider,
) *config {
	g11y.NonNil(h)

	cp := c.open()
	cp.header = cp.header.shallow()
	cp.header.overwriters = maps.Clone(h)
	cp.seal()

	return cp
}

func (c *config) WithBearerToken(
	bt string,
) Config {
	return c.withBearerToken(bt)
}

func (c *config) withBearerToken(
	bt string,
) *config {
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

func (c *config) WithBearerProvider(
	fn HeaderProvider,
) Config {
	return c.withBearerProvider(fn)
}

func (c *config) withBearerProvider(
	fn HeaderProvider,
) *config {
	g11y.NonNil(fn)

	fn = func(ctx context.Context, config Config) string {
		return withBearerPrefix(fn(ctx, config))
	}

	cp := c.open()
	cp.header = cp.header.shallow()
	cp.header.overwriters = maps.Clone(cp.header.overwriters)
	cp.header.overwriters[headers.Authorization] = fn
	cp.seal()

	return cp
}

func (c *config) WithoutBearerProvider() Config {
	return c.withoutBearerProvider()
}

func (c *config) withoutBearerProvider() *config {
	_, ok := c.header.overwriters[headers.Authorization]
	if !ok {
		return c
	}

	cp := c.open()
	cp.header = cp.header.shallow()
	cp.header.overwriters = maps.Clone(cp.header.overwriters)
	delete(cp.header.overwriters, headers.Authorization)
	cp.seal()

	return cp
}

func (c *config) WithExpectingStatusCode(
	code int,
) Config {
	return c.withExpectingStatusCode(code)
}

func (c *config) withExpectingStatusCode(
	code int,
) *config {
	if code < 1 {
		panic("invalid http status code: " + strconv.Itoa(code))
	}

	if c.resp.expectStatusCode == code {
		return c
	}

	cp := c.open()
	cp.resp = cp.resp.shallow()
	cp.resp.expectStatusCode = code
	cp.seal()

	return cp
}

func (c *config) WithoutExpectingStatusCode() Config {
	return c.withoutExpectingStatusCode()
}

func (c *config) withoutExpectingStatusCode() *config {
	if c.resp.expectStatusCode == 0 {
		return c
	}

	cp := c.open()
	cp.resp = cp.resp.shallow()
	cp.resp.expectStatusCode = 0
	cp.seal()

	return cp
}

func (c *config) ExpectingStatusCode() int {
	return c.resp.expectStatusCode
}

func (c *config) Endpoint() string {
	return c.ensure().http.endpoint
}

func (c *config) WithEndpoint(
	e string,
) Config {
	return c.withEndpoint(e)
}

func (c *config) withEndpoint(
	e string,
) *config {
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

func (c *config) WithoutEndpoint() Config {
	return c.withoutEndpoint()
}

func (c *config) withoutEndpoint() *config {
	if c.http.endpoint == "" {
		return c
	}

	cp := c.open()
	cp.http = cp.http.shallow()
	cp.http.endpoint = ""
	cp.seal()

	return cp
}

func (c *config) PathPrefix() string {
	return c.ensure().http.pathPrefix
}

func (c *config) WithPathPrefix(
	p string,
) Config {
	return c.withPathPrefix(p)
}

func (c *config) withPathPrefix(
	p string,
) *config {
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

func (c *config) WithoutPathPrefix() Config {
	return c.withoutPathPrefix()
}

func (c *config) withoutPathPrefix() *config {
	if c.http.pathPrefix == "" {
		return c
	}

	cp := c.open()
	cp.http = cp.http.shallow()
	cp.http.pathPrefix = ""
	cp.seal()

	return cp
}

func (c *config) AndPathPrefix(
	p string,
) Config {
	return c.andPathPrefix(p)
}

func (c *config) andPathPrefix(
	p string,
) *config {
	if p == "" || strings.TrimSpace(p) == "" {
		panic("empty path prefix")
	}

	cp := c.withPathPrefix(Join(c.http.pathPrefix, p))
	cp.seal()

	return cp
}

func (c *config) Timeout() time.Duration {
	return c.ensure().http.timeout
}

func (c *config) WithTimeout(
	t time.Duration,
) Config {
	return c.withTimeout(t)
}

func (c *config) withTimeout(
	t time.Duration,
) *config {
	if c.http.timeout == t {
		return c
	}

	cp := c.open()
	cp.http = cp.http.shallow()
	cp.http.timeout = t
	cp.seal()

	return cp
}

func (c *config) WithTransport(
	t http.RoundTripper,
) Config {
	return c.withTransport(t)
}

func (c *config) withTransport(
	t http.RoundTripper,
) *config {
	g11y.NonNil(t)

	cp := c.open()
	cp.base = t
	cp.seal()

	return cp
}

func (c *config) WithTraced() Config {
	return c.withTraced()
}

func (c *config) withTraced() *config {
	return c.setTraced(true)
}

func (c *config) WithoutTraced() Config {
	return c.withoutTraced()
}

func (c *config) withoutTraced() *config {
	return c.setTraced(false)
}

func (c *config) SetTraced(
	b bool,
) Config {
	return c.setTraced(b)
}

func (c *config) setTraced(
	b bool,
) *config {
	if c.otel.enabled == b {
		return c
	}

	cp := c.open()
	cp.otel = cp.otel.shallow()
	cp.otel.enabled = b
	cp.seal()

	return cp
}

func (c *config) TraceOptions() []otelhttp.Option {
	return slices.Clone(c.ensure().otel.options)
}

func (c *config) WithTraceOptions(
	options ...otelhttp.Option,
) Config {
	return c.withTraceOptions(options...)
}

func (c *config) withTraceOptions(
	options ...otelhttp.Option,
) *config {
	if slices.Equal(c.otel.options, options) {
		return c
	}

	cp := c.open()
	cp.otel = cp.otel.shallow()
	cp.otel.options = options
	cp.seal()

	return cp
}

func (c *config) IsTraced() bool {
	return c.otel.enabled
}

// =============================================================================.

func (c *config) IsRetryLog() bool {
	return c.ensure().retry.logged
}

func (c *config) WithRetryLogged() Config {
	return c.withLogReties()
}

func (c *config) withLogReties() *config {
	return c.setLogReties(true)
}

func (c *config) WithoutRetryLogged() Config {
	return c.withoutLogReties()
}

func (c *config) withoutLogReties() *config {
	return c.setLogReties(false)
}

func (c *config) SetRetryLogged(b bool) Config {
	return c.setLogReties(b)
}

func (c *config) setLogReties(b bool) *config {
	if c.retry.logged == b {
		return c
	}

	cp := c.open()
	cp.retry = cp.retry.shallow()
	cp.retry.logged = b
	cp.seal()

	return cp
}

func (c *config) RetryMax() uint {
	return c.ensure().retry.maxRetries
}

func (c *config) WithMaxRetries(r uint) Config {
	return c.withMaxRetries(r)
}

func (c *config) withMaxRetries(r uint) *config {
	if c.retry.maxRetries == r {
		return c
	}

	cp := c.open()
	cp.retry = cp.retry.shallow()
	cp.retry.maxRetries = r
	cp.seal()

	return cp
}

func (c *config) RetryStatusCodes() []int {
	return slices.Clone(c.ensure().retry.retryIfStatuses)
}

func (c *config) WithRetryStatusCodes(sc ...int) Config {
	return c.withRetriedStatusCodes(sc...)
}

func (c *config) withRetriedStatusCodes(sc ...int) *config {
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

func (c *config) WithRetryIf(fn RetryIfFn) Config {
	return c.withRetryIf(fn)
}

func (c *config) withRetryIf(fn RetryIfFn) *config {
	g11y.NonNil(fn)

	cp := c.open()
	cp.retry = cp.retry.shallow()
	cp.retry.retryIf = fn
	cp.seal()

	return cp
}

func (c *config) WithoutRetryIf() Config {
	return c.withoutRetryIf()
}

func (c *config) withoutRetryIf() *config {
	cp := c.open()
	cp.retry = cp.retry.shallow()
	cp.retry.retryIf = nil
	cp.seal()

	return cp
}

func (c *config) RetryIf() RetryIfFn {
	return c.ensure().retry.retryIf
}

func (c *config) RetryBackoffDuration() time.Duration {
	return c.ensure().retry.backoffDuration
}

func (c *config) WithRetryBackoffDuration(d time.Duration) Config {
	return c.withRetryBackoffDuration(d)
}

func (c *config) withRetryBackoffDuration(d time.Duration) *config {
	if d == c.retry.backoffDuration {
		return c
	}

	cp := c.open()
	cp.retry = cp.retry.shallow()
	cp.retry.backoffDuration = d
	cp.seal()

	return cp
}
