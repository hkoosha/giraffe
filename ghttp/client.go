package ghttp

import (
	"context"
	"io"
	"net/http"

	"github.com/hkoosha/giraffe/g11y"
)

type GClient struct {
	cfg *Config
	hc  *http.Client
}

func NewClient(
	cfg *Config,
) *GClient {
	g11y.NonNil(cfg)
	cfg.ensure()

	return &GClient{
		cfg: cfg,
		hc:  cfg.Std(),
	}
}

func (g *GClient) Std() *http.Client {
	return g.cfg.Std()
}

func (g *GClient) Cfg() *Config {
	return g.cfg
}

func (g *GClient) Get(
	ctx context.Context,
	opts ...GetOptions,
) (*http.Response, error) {
	url := Join(g.cfg.endpoint, g.cfg.pathPrefix)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err2 := opt.apply(req); err2 != nil {
			return nil, err2
		}
	}

	resp, err := g.hc.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *GClient) Post(
	ctx context.Context,
	body io.Reader,
	opts ...PostOptions,
) (*http.Response, error) {
	url := Join(g.cfg.endpoint, g.cfg.pathPrefix)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, body)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err2 := opt.apply(req); err2 != nil {
			return nil, err2
		}
	}

	resp, err := g.hc.Do(req)
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
