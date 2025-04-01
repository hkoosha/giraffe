package ghttp

import (
	"net/http"
)

type giraffeRoundTripper struct {
	cfg *Config
}

func (u giraffeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())

	for k, v := range u.cfg.headerOverwrites {
		req.Header.Set(k, v)
	}

	return u.cfg.baseTransport.RoundTrip(req)
}
