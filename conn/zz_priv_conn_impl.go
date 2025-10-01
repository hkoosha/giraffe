package conn

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/hkoosha/giraffe/conn/internal"
	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/zebra/serdes"
)

type noBodyT struct{}

func (noBodyT) Read([]byte) (int, error)         { return 0, io.EOF }
func (noBodyT) Close() error                     { return nil }
func (noBodyT) WriteTo(io.Writer) (int64, error) { return 0, nil }

var nobody = noBodyT{}

// ============================================================================.

func newConn[R any](
	cfg *config,
	serde serdes.Serde[R],
) *connImpl[R] {
	g11y.NonNil(cfg, serde)
	cfg.Ensure()

	var r R
	return &connImpl[R]{
		Sealer: internal.Sealer{},
		cfg:    cfg,
		std:    cfg.Std(),
		serde:  serde,
		rErr:   r,
	}
}

type connImpl[R any] struct {
	internal.Sealer

	cfg   *config
	std   *http.Client
	serde serdes.Serde[R]
	rErr  R
}

func (c *connImpl[R]) Std() *http.Client {
	return c.cfg.Std()
}

func (c *connImpl[R]) Raw() Conn[[]byte] {
	return newConn[[]byte](c.cfg, serdes.Bytes())
}

func (c *connImpl[R]) Cfg() Config {
	return c.cfg
}

func (c *connImpl[R]) Patch(
	ctx context.Context,
	body any,
	path ...string,
) (R, error) {
	const m = http.MethodPatch
	_, r, err := c.call(ctx, m, body, path)
	return r, err
}

func (c *connImpl[R]) Put(
	ctx context.Context,
	body any,
	path ...string,
) (R, error) {
	const m = http.MethodPut
	_, r, err := c.call(ctx, m, body, path)
	return r, err
}

func (c *connImpl[R]) Post(
	ctx context.Context,
	body any,
	path ...string,
) (R, error) {
	const m = http.MethodPost
	_, r, err := c.call(ctx, m, body, path)
	return r, err
}

func (c *connImpl[R]) PostForHeaders(
	ctx context.Context,
	body any,
	path ...string,
) (http.Header, error) {
	const m = http.MethodPost
	resp, _, err := c.call(ctx, m, body, path)
	if err != nil {
		return nil, err
	}
	return resp.Header, nil
}

func (c *connImpl[R]) Get(
	ctx context.Context,
	path ...string,
) (R, error) {
	const m = http.MethodGet
	_, r, err := c.call(ctx, m, nobody, path)
	return r, err
}

func (c *connImpl[R]) GetForHeaders(
	ctx context.Context,
	path ...string,
) (http.Header, error) {
	const m = http.MethodGet
	resp, _, err := c.call(ctx, m, nobody, path)
	if err != nil {
		return nil, err
	}
	return resp.Header, nil
}

func (c *connImpl[R]) Delete(
	ctx context.Context,
	path ...string,
) (R, error) {
	const m = http.MethodDelete
	_, r, err := c.call(ctx, m, nobody, path)
	return r, err
}

func (c *connImpl[R]) Call(
	ctx context.Context,
	method string,
	body any,
	path ...string,
) (R, error) {
	_, call, err := c.call(ctx, method, body, path)
	return call, err
}

func (c *connImpl[R]) call(
	ctx context.Context,
	method string,
	body any,
	path []string,
) (*http.Response, R, error) {
	var b io.Reader

	switch cast := body.(type) {
	case nil:
		b = nobody

	case io.Reader:
		b = cast

	case []byte:
		b = bytes.NewReader(cast)

	default:
		serialized, err := c.cfg.serde_.Write(cast)
		if err != nil {
			return nil, c.rErr, err
		}
		b = bytes.NewReader(serialized)
	}

	resp, err := c.callRaw(ctx, method, b, path)
	if err != nil {
		return nil, c.rErr, err
	}

	if resp.Body == nil {
		rd, err := c.serde.Read([]byte{})
		return resp, rd, err
	}

	defer resp.Body.Close()

	u, err := c.serde.ReadFrom(resp.Body)
	if err != nil {
		return nil, c.rErr, err
	}

	return resp, u, nil
}

func (c *connImpl[R]) callRaw(
	ctx context.Context,
	method string,
	body io.Reader,
	path []string,
) (*http.Response, error) {
	if body == http.NoBody {
		body = nil
	}

	req, err := http.NewRequestWithContext(ctx, method, join(path), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.std.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
