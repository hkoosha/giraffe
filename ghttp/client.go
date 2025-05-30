package ghttp

import (
	"context"
	"io"
	"net/http"

	"github.com/hkoosha/giraffe/g11y"
)

type (
	HeaderFilter = func(ctx context.Context, h, v string) bool
	HeaderFn     = func(context.Context, *Config) string
)

type Conn struct {
	cfg *Config
	hc  *http.Client
}

func NewClient(
	cfg *Config,
) *Conn {
	g11y.NonNil(cfg)
	cfg.Ensure()

	return &Conn{
		cfg: cfg,
		hc:  cfg.Std(),
	}
}

func (c *Conn) Std() *http.Client {
	return c.cfg.Std()
}

func (c *Conn) Cfg() *Config {
	return c.cfg
}

func (c *Conn) Get(
	ctx context.Context,
	opts ...GetOptions,
) (*http.Response, error) {
	url := Join(c.cfg.Endpoint(), c.cfg.PathPrefix())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err2 := opt.apply(req); err2 != nil {
			return nil, err2
		}
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Conn) Post(
	ctx context.Context,
	body io.Reader,
	opts ...PostOptions,
) (*http.Response, error) {
	url := Join(c.cfg.Endpoint(), c.cfg.PathPrefix())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, body)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err2 := opt.apply(req); err2 != nil {
			return nil, err2
		}
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type GetOptions interface {
	apply(req *http.Request) error

	seal(seal)
}

type PostOptions interface {
	apply(req *http.Request) error

	seal(seal)
}
