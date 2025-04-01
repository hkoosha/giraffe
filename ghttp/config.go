package ghttp

import (
	"context"
	"maps"
	"net/http"
	"strings"
	"time"

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/ghttp/headers"
	"github.com/hkoosha/giraffe/glog"
	. "github.com/hkoosha/giraffe/internal/dot"
	"github.com/hkoosha/giraffe/zebra/z"
)

const (
	UserAgent = "Giraffe/1.0"
	Timeout   = 5 * time.Second
)

var defaultHeaders = map[string]string{
	headers.UserAgent: UserAgent,
}

func NewConfig(
	lg glog.GLog,
) *Config {
	cfg := &Config{
		lg:                 lg,
		isLogged:           true,
		isPlainLog:         false,
		lgFilteredHeaders:  map[string]z.NA{},
		lgMaskedHeaders:    map[string]z.NA{},
		isOtel:             true,
		headerOverwrites:   map[string]string{},
		bearerProvider:     nil,
		isLogReties:        false,
		retriedStatusCodes: []int{},
		isExpecting2xx:     true,
		endpoint:           "",
		pathPrefix:         "",
		timeout:            Timeout,
		baseTransport:      http.DefaultTransport,
		transport:          giraffeRoundTripper{cfg: nil},
		sealed:             false,
	}

	cfg.transport.cfg = cfg

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
	baseTransport      http.RoundTripper
	lg                 glog.GLog
	bearerProvider     func(ctx context.Context) string
	transport          giraffeRoundTripper
	lgFilteredHeaders  z.Set[string]
	lgMaskedHeaders    z.Set[string]
	headerOverwrites   map[string]string
	pathPrefix         string
	endpoint           string
	retriedStatusCodes []int
	timeout            time.Duration
	isLogReties        bool
	isExpecting2xx     bool
	isOtel             bool
	isPlainLog         bool
	isLogged           bool
	sealed             bool
}

func (c *Config) ensure() *Config {
	if !c.sealed {
		panic(EF("invalid config, did you use constructor to create one?"))
	}

	return c
}

func (c *Config) open() *Config {
	c.ensure()

	return &Config{
		lg:                 c.lg,
		isLogged:           c.isLogged,
		isPlainLog:         c.isPlainLog,
		lgFilteredHeaders:  c.lgFilteredHeaders,
		lgMaskedHeaders:    c.lgMaskedHeaders,
		isOtel:             c.isOtel,
		headerOverwrites:   c.headerOverwrites,
		bearerProvider:     c.bearerProvider,
		isLogReties:        c.isLogReties,
		retriedStatusCodes: c.retriedStatusCodes,
		isExpecting2xx:     c.isExpecting2xx,
		endpoint:           c.endpoint,
		pathPrefix:         c.pathPrefix,
		timeout:            c.timeout,
		baseTransport:      c.baseTransport,
		transport:          giraffeRoundTripper{cfg: nil},
		sealed:             false,
	}
}

func (c *Config) seal() {
	if c.sealed {
		return
	}

	for k := range c.lgMaskedHeaders {
		if _, ok := c.lgFilteredHeaders[k]; ok {
			panic(EF("%s", "header cannot be both filtered and masked: "+k))
		}
	}

	for k := range c.lgFilteredHeaders {
		if _, ok := c.lgMaskedHeaders[k]; ok {
			panic(EF("%s", "header cannot be both filtered and masked: "+k))
		}
	}

	c.transport = giraffeRoundTripper{cfg: c}

	c.sealed = true
}

func (c *Config) Std() *http.Client {
	return c.ensure().std0()
}

func (c *Config) std0() *http.Client {
	return &http.Client{
		Transport:     c.transport,
		Timeout:       c.timeout,
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
	return c.ensure().isLogged
}

func (c *Config) WithLogged() *Config {
	return c.SetLogged(true)
}

func (c *Config) WithoutLogged() *Config {
	return c.SetLogged(false)
}

func (c *Config) SetLogged(b bool) *Config {
	if c.sealed && c.isLogged == b {
		return c
	}

	cp := c.open()
	cp.isLogged = b
	cp.seal()

	return cp
}

func (c *Config) IsPlainLog() bool {
	return c.ensure().isPlainLog
}

func (c *Config) WithPlainLog() *Config {
	return c.SetPlainLog(true)
}

func (c *Config) WithoutPlainLog() *Config {
	return c.SetPlainLog(false)
}

func (c *Config) SetPlainLog(b bool) *Config {
	if c.sealed && c.isPlainLog == b {
		return c
	}

	cp := c.open()
	cp.isPlainLog = b
	cp.seal()

	return cp
}

func (c *Config) LgFilteredHeaders() z.Set[string] {
	return maps.Clone(c.ensure().lgFilteredHeaders)
}

func (c *Config) WithLgFilteredHeaders(h z.Set[string]) *Config {
	g11y.NonNil(h)
	if c.sealed && z.MapEq(c.lgFilteredHeaders, h) {
		return c
	}

	cp := c.open()
	cp.lgFilteredHeaders = maps.Clone(h)
	cp.seal()

	return cp
}

func (c *Config) LgMaskedHeaders() z.Set[string] {
	return maps.Clone(c.ensure().lgMaskedHeaders)
}

func (c *Config) WithLgMaskedHeaders(h z.Set[string]) *Config {
	g11y.NonNil(h)
	if c.sealed && z.MapEq(c.lgMaskedHeaders, h) {
		return c
	}

	cp := c.open()
	cp.lgMaskedHeaders = maps.Clone(h)
	cp.seal()

	return cp
}

func (c *Config) IsOtel() bool {
	return c.ensure().isOtel
}

func (c *Config) WithOtel() *Config {
	return c.SetOtel(true)
}

func (c *Config) WithoutOtel() *Config {
	return c.SetOtel(false)
}

func (c *Config) SetOtel(b bool) *Config {
	if c.sealed && c.isOtel == b {
		return c
	}

	cp := c.open()
	cp.isOtel = b
	cp.seal()

	return cp
}

func (c *Config) HeaderOverwrites() map[string]string {
	return maps.Clone(c.ensure().headerOverwrites)
}

func (c *Config) WithHeaderOverwrites(
	withDefaults bool,
	h map[string]string,
) *Config {
	g11y.NonNil(h)
	if withDefaults {
		h = z.UnionLeft(h, defaultHeaders)
	}

	if c.sealed && z.MapEq(c.headerOverwrites, h) {
		return c
	}

	h = maps.Clone(h)

	cp := c.open()
	cp.headerOverwrites = maps.Clone(h)
	cp.seal()

	return cp
}

func (c *Config) WithBearerToken(bt string) *Config {
	if strings.TrimSpace(bt) == "" {
		panic("empty bearer token")
	}

	fn := func(ctx context.Context) string {
		return bt
	}

	return c.withBearerProvider0(fn)
}

func (c *Config) WithBearerProvider(fn func(context.Context) string) *Config {
	g11y.NonNil(fn)

	return c.withBearerProvider0(fn)
}

func (c *Config) WithoutBearerProvider() *Config {
	return c.withBearerProvider0(nil)
}

func (c *Config) withBearerProvider0(fn func(context.Context) string) *Config {
	if fn == nil && c.bearerProvider == nil {
		return c
	}

	cp := c.open()
	cp.bearerProvider = fn
	cp.seal()

	return cp
}

func (c *Config) IsLogReties() bool {
	c.ensure()
	panic("unimplemented")
}

func (c *Config) WithLogReties() *Config {
	return c.SetLogReties(true)
}

func (c *Config) WithoutLogReties() *Config {
	return c.SetLogReties(false)
}

func (c *Config) SetLogReties(_ bool) *Config {
	panic("unimplemented")
}

func (c *Config) MaxRetries() uint {
	c.ensure()
	panic("unimplemented")
}

func (c *Config) WithMaxRetries(r uint) *Config {
	c.ensure()

	if r > 50 {
		panic(EF("max retries too large: %d", r))
	}

	panic("unimplemented")
}

func (c *Config) RetriedStatusCodes() []int {
	c.ensure()
	panic("unimplemented")
}

func (c *Config) WithRetriedStatusCodes(sc ...int) *Config {
	g11y.NonNil(sc)
	c.ensure()
	panic("unimplemented")
}

func (c *Config) IsExpecting2xx() bool {
	return c.ensure().isExpecting2xx
}

func (c *Config) WithExpecting2xx() *Config {
	return c.SetExpecting2xx(true)
}

func (c *Config) WithoutExpecting2xx() *Config {
	return c.SetExpecting2xx(false)
}

func (c *Config) SetExpecting2xx(b bool) *Config {
	if c.sealed && c.isExpecting2xx == b {
		return c
	}

	cp := c.open()
	cp.isExpecting2xx = b
	cp.seal()

	return c
}

func (c *Config) Endpoint() string {
	return c.ensure().endpoint
}

func (c *Config) WithEndpoint(e string) *Config {
	if c.sealed && c.endpoint == e {
		return c
	}

	cp := c.open()
	cp.endpoint = e
	cp.seal()

	return cp
}

func (c *Config) PathPrefix() string {
	return c.ensure().pathPrefix
}

func (c *Config) WithPathPrefix(p string) *Config {
	// TODO validate p.
	if p == "" || strings.TrimSpace(p) == "" {
		panic("empty path prefix")
	}

	if c.sealed && c.pathPrefix == p {
		return c
	}

	cp := c.open()
	cp.pathPrefix = p
	cp.seal()

	return cp
}

func (c *Config) AndPathPrefix(p string) *Config {
	if p == "" || strings.TrimSpace(p) == "" {
		panic("empty path prefix")
	}

	return c.WithPathPrefix(Join(c.pathPrefix, p))
}

func (c *Config) Timeout() time.Duration {
	return c.ensure().timeout
}

func (c *Config) WithTimeout(t time.Duration) *Config {
	if c.sealed && c.timeout == t {
		return c
	}

	cp := c.open()
	cp.timeout = t
	cp.seal()

	return cp
}

func (c *Config) WithTransport(transport http.RoundTripper) *Config {
	g11y.NonNil(transport)
	cp := c.open()
	cp.baseTransport = transport

	return cp
}

type seal struct{}
