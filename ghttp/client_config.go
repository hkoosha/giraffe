package ghttp

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/hkoosha/giraffe/glog"
)

func (c *Conn) WithLg(lg glog.GLog) *Conn {
	return NewClient(c.cfg.WithLg(lg))
}

func (c *Conn) WithLogged() *Conn {
	return NewClient(c.cfg.WithLogged())
}

func (c *Conn) WithoutLogged() *Conn {
	return NewClient(c.cfg.WithoutLogged())
}

func (c *Conn) WithPlainLog() *Conn {
	return NewClient(c.cfg.WithPlainLog())
}

func (c *Conn) WithoutPlainLog() *Conn {
	return NewClient(c.cfg.WithoutPlainLog())
}

func (c *Conn) WithLgFilteredHeaders(f HeaderFilter) *Conn {
	return NewClient(c.cfg.WithLgFilteredHeaders(f))
}

func (c *Conn) WithLgMaskedHeaders(f HeaderFilter) *Conn {
	return NewClient(c.cfg.WithLgMaskedHeaders(f))
}

func (c *Conn) WithHeaderOverwrites(
	withDefaults bool,
	h map[string]string,
) *Conn {
	return NewClient(c.cfg.WithHeaderOverwrites(withDefaults, h))
}

func (c *Conn) WithBearerToken(bt string) *Conn {
	return NewClient(c.cfg.WithBearerToken(bt))
}

func (c *Conn) WithBearerProvider(fn HeaderFn) *Conn {
	return NewClient(c.cfg.WithBearerProvider(fn))
}

func (c *Conn) WithoutBearerProvider() *Conn {
	return NewClient(c.cfg.WithoutBearerProvider())
}

func (c *Conn) WithExpecting2xx() *Conn {
	return NewClient(c.cfg.WithExpecting2xx())
}

func (c *Conn) WithoutExpecting2xx() *Conn {
	return NewClient(c.cfg.WithoutExpecting2xx())
}

func (c *Conn) WithEndpoint(e string) *Conn {
	return NewClient(c.cfg.WithEndpoint(e))
}

func (c *Conn) WithoutEndpoint() *Conn {
	return NewClient(c.cfg.WithoutEndpoint())
}

func (c *Conn) WithPathPrefix(p string) *Conn {
	return NewClient(c.cfg.WithPathPrefix(p))
}

func (c *Conn) WithoutPathPrefix() *Conn {
	return NewClient(c.cfg.WithoutPathPrefix())
}

func (c *Conn) AndPathPrefix(p string) *Conn {
	return NewClient(c.cfg.AndPathPrefix(p))
}

func (c *Conn) WithTimeout(t time.Duration) *Conn {
	return NewClient(c.cfg.WithTimeout(t))
}

func (c *Conn) WithTransport(transport http.RoundTripper) *Conn {
	return NewClient(c.cfg.WithTransport(transport))
}

func (c *Conn) WithTraced() *Conn {
	return NewClient(c.cfg.WithTraced())
}

func (c *Conn) WithoutTraced() *Conn {
	return NewClient(c.cfg.WithoutTraced())
}

func (c *Conn) WithTraceOptions(
	options ...otelhttp.Option,
) *Conn {
	return NewClient(c.cfg.WithTraceOptions(options...))
}

// =============================================================================.

func (c *Conn) WithLogReties() *Conn {
	return NewClient(c.cfg.WithLogReties())
}

func (c *Conn) WithoutLogReties() *Conn {
	return NewClient(c.cfg.WithoutLogReties())
}

func (c *Conn) WithMaxRetries(r uint) *Conn {
	return NewClient(c.cfg.WithMaxRetries(r))
}

func (c *Conn) WithRetriedStatusCodes(sc ...int) *Conn {
	return NewClient(c.cfg.WithRetriedStatusCodes(sc...))
}
