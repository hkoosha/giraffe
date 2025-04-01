package ghttp

import (
	"context"
	"net/http"
	"time"

	"github.com/hkoosha/giraffe/glog"
	"github.com/hkoosha/giraffe/zebra/z"
)

func (g *GClient) WithLg(
	lg glog.GLog,
) *GClient {
	return NewClient(g.cfg.WithLg(lg))
}

func (g *GClient) WithLogged() *GClient {
	return NewClient(g.cfg.WithLogged())
}

func (g *GClient) WithoutLogged() *GClient {
	return NewClient(g.cfg.WithoutLogged())
}

func (g *GClient) WithPlainLog() *GClient {
	return NewClient(g.cfg.WithPlainLog())
}

func (g *GClient) WithoutPlainLog() *GClient {
	return NewClient(g.cfg.WithoutPlainLog())
}

func (g *GClient) WithLgFilteredHeaders(
	h z.Set[string],
) *GClient {
	return NewClient(g.cfg.WithLgFilteredHeaders(h))
}

func (g *GClient) WithLgMaskedHeaders(
	h z.Set[string],
) *GClient {
	return NewClient(g.cfg.WithLgMaskedHeaders(h))
}

func (g *GClient) WithOtel() *GClient {
	return NewClient(g.cfg.WithOtel())
}

func (g *GClient) WithoutOtel() *GClient {
	return NewClient(g.cfg.WithoutOtel())
}

func (g *GClient) WithHeaderOverwrites(
	withDefaults bool,
	h map[string]string,
) *GClient {
	return NewClient(g.cfg.WithHeaderOverwrites(withDefaults, h))
}

func (g *GClient) WithBearerToken(bt string) *GClient {
	return NewClient(g.cfg.WithBearerToken(bt))
}

func (g *GClient) WithBearerProvider(fn func(context.Context) string) *GClient {
	return NewClient(g.cfg.WithBearerProvider(fn))
}

func (g *GClient) WithoutBearerProvider() *GClient {
	return NewClient(g.cfg.WithoutBearerProvider())
}

func (g *GClient) WithLogReties() *GClient {
	return NewClient(g.cfg.WithLogReties())
}

func (g *GClient) WithoutLogReties() *GClient {
	return NewClient(g.cfg.WithoutLogReties())
}

func (g *GClient) SetLogReties(b bool) *GClient {
	return NewClient(g.cfg.SetLogReties(b))
}

func (g *GClient) WithMaxRetries(r uint) *GClient {
	return NewClient(g.cfg.WithMaxRetries(r))
}

func (g *GClient) WithRetriedStatusCodes(sc ...int) *GClient {
	return NewClient(g.cfg.WithRetriedStatusCodes(sc...))
}

func (g *GClient) WithExpecting2xx() *GClient {
	return NewClient(g.cfg.WithExpecting2xx())
}

func (g *GClient) WithoutExpecting2xx() *GClient {
	return NewClient(g.cfg.WithoutExpecting2xx())
}

func (g *GClient) SetExpecting2xx(b bool) *GClient {
	return NewClient(g.cfg.SetExpecting2xx(b))
}

func (g *GClient) WithEndpoint(e string) *GClient {
	return NewClient(g.cfg.WithEndpoint(e))
}

func (g *GClient) WithPathPrefix(p string) *GClient {
	return NewClient(g.cfg.WithPathPrefix(p))
}

func (g *GClient) AndPathPrefix(p string) *GClient {
	return NewClient(g.cfg.AndPathPrefix(p))
}

func (g *GClient) WithTimeout(t time.Duration) *GClient {
	return NewClient(g.cfg.WithTimeout(t))
}

func (g *GClient) WithTransport(transport http.RoundTripper) *GClient {
	return NewClient(g.cfg.WithTransport(transport))
}
