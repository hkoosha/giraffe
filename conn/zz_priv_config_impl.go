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

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn/headers"
	"github.com/hkoosha/giraffe/conn/internal"
	"github.com/hkoosha/giraffe/core/serdes"
	"github.com/hkoosha/giraffe/core/t11y"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/glog"
	"github.com/hkoosha/giraffe/zebra/z"
)

// TODO call ensure on With...()

func cfgOf(
	cfg Config,
) *config {
	cfg.Ensure()

	if cast, ok := cfg.(*config); ok {
		return cast
	}

	cast := newConfig(cfg.Lg(), cfg.Timeout()).
		withTxSerde(cfg.TxSerde()).
		withRxSerde(cfg.RxSerde()).
		setPlainLog(cfg.IsPlainLog()).
		setLogged(cfg.IsLogged()).
		setLogReties(cfg.IsRetryLog()).
		withMaxRetries(cfg.RetryMax()).
		setTraced(cfg.IsTraced()).
		withRetryBackoffDuration(cfg.RetryBackoffDuration())

	if v := cfg.ExpectingStatusCodes(); len(v) > 0 {
		cast = cast.withExpectingStatusCodes(v)
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
			panic(EF("invalid config, unsupported header overwriter: %s", header))
		}
	}
	cast = cast.withHeaderOverwriters(cfg.HeaderOverwriters())

	return cast
}

func newConfig(
	lg glog.Lg,
	timeout time.Duration,
) *config {
	cfg := &config{
		Sealer:   internal.Sealer{},
		sealed:   false,
		base:     defaultTransport,
		lg:       lg,
		rt:       nil,
		resp:     mkResponseConfig(),
		http:     mkHttpConfig(timeout),
		header:   mkHeaderConfig(),
		log:      mkLogConfig(),
		retry:    mkRetryConfig(),
		otel:     mkOtelConfig(),
		txSerde_: serdes.Bytes(),
		rxSerde_: serdes.Bytes(),
	}

	cfg.seal()

	return cfg
}

type config struct {
	internal.Sealer

	txSerde_ any
	rxSerde_ any

	lg     glog.Lg
	base   http.RoundTripper
	rt     http.RoundTripper
	resp   *respConfig
	http   *httpConfig
	header *headerConfig
	log    *logConfig
	retry  *retryConfig
	otel   *otelConfig

	sealed bool
}

// =============================================================================.

func (c *config) RxSerde() any {
	return c.rxSerde()
}

func (c *config) rxSerde() any {
	c.ensure()

	return c.rxSerde_
}

func (c *config) TxSerde() any {
	return c.txSerde()
}

func (c *config) txSerde() any {
	c.ensure()

	return c.txSerde_
}

func (c *config) WithSerdes(
	tx any,
	rx any,
) Config {
	return c.withSerdes(tx, rx)
}

func (c *config) withSerdes(
	tx any,
	rx any,
) *config {
	return c.withRxSerde(rx).withTxSerde(tx)
}

func (c *config) WithTxSerde(v any) Config {
	return c.withTxSerde(v)
}

func (c *config) withTxSerde(v any) *config {
	if !serdes.IsSerde(v) {
		panic(EF("not a serde: %v", v))
	}

	cp := c.open()
	cp.txSerde_ = v
	cp.seal()

	return cp
}

func (c *config) WithRxSerde(v any) Config {
	return c.withRxSerde(v)
}

func (c *config) withRxSerde(v any) *config {
	if !serdes.IsSerde(v) {
		panic(EF("not a serde: %v", v))
	}

	cp := c.open()
	cp.rxSerde_ = v
	cp.seal()

	return cp
}

func (c *config) WithJsonTxSerde() Config {
	return c.withJsonTxSerde()
}

func (c *config) withJsonTxSerde() *config {
	return c.withTxSerde(serdes.Json[any]())
}

func (c *config) WithJsonRxSerde() Config {
	return c.withJsonTxSerde()
}

func (c *config) withJsonRxSerde() *config {
	return c.withRxSerde(serdes.Json[any]())
}

func (c *config) WithJsonSerde() Config {
	return c.withJsonSerde()
}

func (c *config) withJsonSerde() *config {
	return c.withJsonRxSerde().withJsonTxSerde()
}

func (c *config) WithStringRxSerde() Config {
	return c.withStringRxSerde()
}

func (c *config) withStringRxSerde() Config {
	return c.withRxSerde(serdes.String())
}

func (c *config) WithBytesRxSerde() Config {
	return c.withBytesRxSerde()
}

func (c *config) withBytesRxSerde() *config {
	return c.withRxSerde(serdes.Bytes())
}

func (c *config) WithBytesTxSerde() Config {
	return c.withBytesTxSerde()
}

func (c *config) withBytesTxSerde() *config {
	return c.withTxSerde(serdes.Bytes())
}

func (c *config) WithBytesSerde() Config {
	return c.withBytesSerde()
}

func (c *config) withBytesSerde() *config {
	return c.withBytesRxSerde().withBytesTxSerde()
}

func (c *config) Raw() Raw {
	return c.connection()
}

func (c *config) Datum() Datum {
	return c.datum()
}

func (c *config) connection() Raw {
	c.ensure()

	return newConn[[]byte, []byte](c.withBytesSerde())
}

func (c *config) datum() Datum {
	c.ensure()

	j := c.withSerdes(giraffe.DatumSerde(), giraffe.DatumSerde())
	return newConn[giraffe.Datum, giraffe.Datum](j)
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
	t11y.NonNil(lg)

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
	t11y.NonNil(f)

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
	t11y.NonNil(f)

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

func (c *config) AndHeaders(
	h map[string]string,
) Config {
	return c.withHeaderOverwrites(true, h)
}

func (c *config) AndHeader(
	name string,
	value string,
) Config {
	return c.withHeaderOverwrites(true, map[string]string{
		name: value,
	})
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
	t11y.NonNil(h)

	if withDefaults {
		h = z.UnionLeft(h, defaultHeaders)
	} else {
		h = maps.Clone(h)
	}

	if z.Eq2(c.header.overwrite, h) {
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
	t11y.NonNil(h)

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
		panic(EF("empty bearer token"))
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
	t11y.NonNil(fn)

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

func (c *config) WithExpectingStatusCodes(
	code ...int,
) Config {
	return c.withExpectingStatusCodes(code)
}

func (c *config) withExpectingStatusCodes(
	codes []int,
) *config {
	for _, code := range codes {
		if code < 1 {
			panic(EF("invalid http status code: %d", code))
		}
	}

	if slices.Equal(c.resp.expectStatusCode, codes) {
		return c
	}

	cp := c.open()
	cp.resp = cp.resp.clone()
	cp.resp.expectStatusCode = slices.Clone(codes)
	cp.seal()

	return cp
}

func (c *config) WithoutExpectingStatusCodes() Config {
	return c.withoutExpectingStatusCodes()
}

func (c *config) withoutExpectingStatusCodes() *config {
	if c.resp.expectStatusCode == nil {
		return c
	}

	cp := c.open()
	cp.resp = cp.resp.clone()
	cp.resp.expectStatusCode = nil
	cp.seal()

	return cp
}

func (c *config) ExpectingStatusCodes() []int {
	return slices.Clone(c.resp.expectStatusCode)
}

func (c *config) Endpoints() map[string]string {
	return maps.Clone(c.ensure().http.endpointsByName)
}

func (c *config) AndEndpoint(
	name string,
	addr string,
) Config {
	return c.andEndpoints(name, addr)
}

func (c *config) andEndpoints(
	name string,
	addr string,
) *config {
	e := maps.Clone(c.http.endpointsByName)
	e[name] = addr
	return c.withEndpoints(e)
}

func (c *config) WithEndpoints(
	e map[string]string,
) Config {
	return c.withEndpoints(e)
}

func (c *config) withEndpoints(
	e map[string]string,
) *config {
	if maps.Equal(c.http.endpointsByName, e) {
		return c
	}

	for name, addr := range e {
		if !endpointNameRe.MatchString(name) {
			panic(EF("invalid endpoint name: %s", name))
		}
		if !endpointAddrRe.MatchString(addr) {
			panic(EF("invalid endpoint address: %s", addr))
		}
		matches := endpointAddrRe.FindStringSubmatch(addr)
		groups := make(map[string]string)
		for i, n := range endpointAddrReNames {
			if i != 0 && n != "" {
				groups[name] = matches[i]
			}
		}

		if strings.Contains(groups["address"], "..") {
			panic(EF("invalid endpoint address: %s", addr))
		}

		if groups["port"] != "" {
			port := M(strconv.Atoi(groups["port"]))
			if port < 1 || 65534 < port {
				panic(EF("invalid endpoint port: %s", addr))
			}
		}
	}

	cp := c.open()
	cp.http = cp.http.shallow()
	cp.http.endpointsByName = maps.Clone(e)
	cp.seal()

	return cp
}

func (c *config) WithoutEndpoints() Config {
	return c.withoutEndpoints()
}

func (c *config) withoutEndpoints() *config {
	if c.http.endpoint == "" {
		return c
	}

	cp := c.open()
	cp.http = cp.http.shallow()
	cp.http.endpoint = ""
	cp.seal()

	return cp
}

func (c *config) Endpoint() string {
	return c.ensure().http.endpoint
}

func (c *config) WithMustEndpointNamed(
	name string,
) Config {
	return M(c.WithEndpointNamed(name))
}

func (c *config) WithEndpointNamed(
	name string,
) (Config, error) {
	return c.withEndpointNamed(name)
}

func (c *config) withEndpointNamed(
	name string,
) (*config, error) {
	if name == "" || strings.TrimSpace(name) == "" {
		panic(EF("empty endpoint name"))
	}

	ep, ok := c.http.endpointsByName[name]
	if !ok {
		return nil, &MissingEndpointError{Endpoint: name}
	}

	return c.withEndpoint(ep), nil
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
		panic(EF("empty endpoint"))
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

func (c *config) Method() string {
	return c.ensure().http.defaultMethod
}

func (c *config) WithMethod(
	p string,
) Config {
	return c.withMethod(p)
}

func (c *config) withMethod(
	v string,
) *config {
	// TODO validate p.
	if v == "" || strings.TrimSpace(v) == "" {
		panic(EF("empty method"))
	}

	if c.http.defaultMethod == v {
		return c
	}

	cp := c.open()
	cp.http = cp.http.shallow()
	cp.http.defaultMethod = v
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
		panic(EF("empty path prefix"))
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
		panic(EF("empty path prefix"))
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
	t11y.NonNil(t)

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
	t11y.NonNil(sc)

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
	t11y.NonNil(fn)

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
