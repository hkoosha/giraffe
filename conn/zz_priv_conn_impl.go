package conn

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/hkoosha/giraffe/conn/internal"
	"github.com/hkoosha/giraffe/internal/reflected"
	"github.com/hkoosha/giraffe/t11y"
	. "github.com/hkoosha/giraffe/t11y/dot"
	"github.com/hkoosha/giraffe/zebra/serdes"
)

func newConn[TX, RX any](
	cfg *config,
) *connImpl[TX, RX] {
	t11y.NonNil(cfg)

	cfg.Ensure()

	return &connImpl[TX, RX]{
		Sealer:  internal.Sealer{},
		cfg:     cfg,
		std:     cfg.Std(),
		rxSerde: serdes.MustCast[RX](cfg.rxSerde()),
		txSerde: serdes.MustCast[TX](cfg.txSerde()),
		rxErr:   reflected.Zero[RX](),
	}
}

type connImpl[TX, RX any] struct {
	internal.Sealer

	cfg *config
	std *http.Client

	rxSerde serdes.Serde[RX]
	txSerde serdes.Serde[TX]
	rxErr   RX
}

func (c *connImpl[TX, RX]) Std() *http.Client {
	return c.cfg.Std()
}

func (c *connImpl[TX, RX]) Raw() Conn[[]byte, []byte] {
	return newConn[[]byte, []byte](
		c.cfg.
			withRxSerde(serdes.Bytes()).
			withTxSerde(serdes.Bytes()),
	)
}

func (c *connImpl[TX, RX]) Cfg() Config {
	return c.cfg
}

func (c *connImpl[TX, RX]) Patch(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, error) {
	const m = http.MethodPatch
	_, r, err := c.call(ctx, m, &body, path)
	return r, err
}

func (c *connImpl[TX, RX]) Put(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, error) {
	const m = http.MethodPut
	_, r, err := c.call(ctx, m, &body, path)
	return r, err
}

func (c *connImpl[TX, RX]) Post(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, error) {
	const m = http.MethodPost
	_, r, err := c.call(ctx, m, &body, path)
	return r, err
}

func (c *connImpl[TX, RX]) PostForHeaders(
	ctx context.Context,
	body TX,
	path ...string,
) (http.Header, error) {
	const m = http.MethodPost
	headers, _, err := c.call(ctx, m, &body, path)
	if err != nil {
		return nil, err
	}
	return headers, nil
}

func (c *connImpl[TX, RX]) Get(
	ctx context.Context,
	path ...string,
) (RX, error) {
	const m = http.MethodGet
	_, r, err := c.call(ctx, m, nil, path)
	return r, err
}

func (c *connImpl[TX, RX]) GetForHeaders(
	ctx context.Context,
	path ...string,
) (http.Header, error) {
	const m = http.MethodGet
	headers, _, err := c.call(ctx, m, nil, path)
	if err != nil {
		return nil, err
	}
	return headers, nil
}

func (c *connImpl[TX, RX]) Delete(
	ctx context.Context,
	path ...string,
) (RX, error) {
	const m = http.MethodDelete
	_, r, err := c.call(ctx, m, nil, path)
	return r, err
}

func (c *connImpl[TX, RX]) Call(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, error) {
	_, call, err := c.call(ctx, c.cfg.http.defaultMethod, &body, path)
	return call, err
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) Headered(
	ctx context.Context,
	body TX,
	path ...string,
) (_ RX, headers map[string]string, _ error) {
	respHeaders, call, err := c.call(ctx, c.cfg.http.defaultMethod, &body, path)
	if err != nil {
		return call, nil, err
	}

	headers0 := make(map[string]string, len(respHeaders))
	for k, v := range respHeaders {
		if len(v) > 0 {
			panic("todo multivalued headers")
		}
		if len(v) == 0 {
			panic("todo: empty headers")
		}

		headers0[k] = v[0]
	}

	return call, headers0, err
}

func (c *connImpl[TX, RX]) call(
	ctx context.Context,
	method string,
	reqBody *TX,
	path []string,
) (http.Header, RX, error) {
	serialized, err := c.txSerde.Write(*reqBody)
	if err != nil {
		return nil, c.rxErr, E(err)
	}

	resp, err := c.callRaw(ctx, method, bytes.NewReader(serialized), path)
	if err != nil {
		return nil, c.rxErr, err
	}

	defer resp.Body.Close()

	u, err := c.rxSerde.StreamFrom(resp.Body)
	if err != nil {
		return nil, c.rxErr, err
	}

	return resp.Header, u, nil
}

func (c *connImpl[TX, RX]) callRaw(
	ctx context.Context,
	method string,
	body io.Reader,
	path []string,
) (*http.Response, error) {
	if body == nil {
		body = nobody
	}

	req, err := http.NewRequestWithContext(ctx, method, join(path), body)
	if err != nil {
		return nil, E(err)
	}

	resp, err := c.std.Do(req)
	if err != nil {
		return nil, E(err)
	}

	return resp, nil
}
