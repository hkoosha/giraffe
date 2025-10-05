package conn

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/hkoosha/giraffe/conn/internal"
	"github.com/hkoosha/giraffe/t11y"
	. "github.com/hkoosha/giraffe/t11y/dot"
	"github.com/hkoosha/giraffe/zebra/serdes"
)

func newConn[RX, TX any](
	cfg *config,
) *connImpl[RX, TX] {
	t11y.NonNil(cfg)

	cfg.Ensure()

	var rx RX
	var tx TX
	return &connImpl[RX, TX]{
		Sealer:  internal.Sealer{},
		cfg:     cfg,
		std:     cfg.Std(),
		rxSerde: serdes.MustCast[RX](cfg.rxSerde()),
		txSerde: serdes.MustCast[TX](cfg.txSerde()),
		rxErr:   rx,
		txErr:   tx,
	}
}

type connImpl[RX, TX any] struct {
	internal.Sealer

	cfg *config
	std *http.Client

	rxSerde serdes.Serde[RX]
	txSerde serdes.Serde[TX]
	rxErr   RX
	txErr   TX
}

func (c *connImpl[RX, TX]) Std() *http.Client {
	return c.cfg.Std()
}

func (c *connImpl[RX, TX]) Raw() Conn[[]byte, []byte] {
	return newConn[[]byte, []byte](
		c.cfg.
			withRxSerde(serdes.Bytes()).
			withTxSerde(serdes.Bytes()),
	)
}

func (c *connImpl[RX, TX]) Cfg() Config {
	return c.cfg
}

func (c *connImpl[RX, TX]) Patch(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, error) {
	const m = http.MethodPatch
	//nolint:bodyclose
	_, r, err := c.call(ctx, m, &body, path)
	return r, err
}

func (c *connImpl[RX, TX]) Put(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, error) {
	const m = http.MethodPut
	//nolint:bodyclose
	_, r, err := c.call(ctx, m, &body, path)
	return r, err
}

func (c *connImpl[RX, TX]) Post(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, error) {
	const m = http.MethodPost
	//nolint:bodyclose
	_, r, err := c.call(ctx, m, &body, path)
	return r, err
}

func (c *connImpl[RX, TX]) PostForHeaders(
	ctx context.Context,
	body TX,
	path ...string,
) (http.Header, error) {
	const m = http.MethodPost
	//nolint:bodyclose
	resp, _, err := c.call(ctx, m, &body, path)
	if err != nil {
		return nil, err
	}
	return resp.Header, nil
}

func (c *connImpl[RX, TX]) Get(
	ctx context.Context,
	path ...string,
) (RX, error) {
	const m = http.MethodGet
	//nolint:bodyclose
	_, r, err := c.call(ctx, m, nil, path)
	return r, err
}

func (c *connImpl[RX, TX]) GetForHeaders(
	ctx context.Context,
	path ...string,
) (http.Header, error) {
	const m = http.MethodGet
	//nolint:bodyclose
	resp, _, err := c.call(ctx, m, nil, path)
	if err != nil {
		return nil, err
	}
	return resp.Header, nil
}

func (c *connImpl[RX, TX]) Delete(
	ctx context.Context,
	path ...string,
) (RX, error) {
	const m = http.MethodDelete
	//nolint:bodyclose
	_, r, err := c.call(ctx, m, nil, path)
	return r, err
}

func (c *connImpl[RX, TX]) Call(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, error) {
	_, call, err := c.call(ctx, c.cfg.http.defaultMethod, &body, path)
	return call, err
}

func (c *connImpl[RX, TX]) CallWithHeaders(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, map[string]string, error) {
	//nolint:bodyclose
	resp, call, err := c.call(ctx, c.cfg.http.defaultMethod, &body, path)
	if err != nil {
		return call, nil, err
	}

	headers := make(map[string]string, len(resp.Header))
	for k, v := range resp.Header {
		if len(v) > 0 {
			panic("todo multivalued headers")
		}
		if len(v) == 0 {
			panic("todo: empty headers")
		}

		headers[k] = v[0]
	}

	return call, headers, err
}

func (c *connImpl[RX, TX]) call(
	ctx context.Context,
	method string,
	reqBody *TX,
	path []string,
) (*http.Response, RX, error) {
	serialized, err := c.txSerde.Write(*reqBody)
	if err != nil {
		return nil, c.rxErr, E(err)
	}

	b := bytes.NewReader(serialized)

	resp, err := c.callRaw(ctx, method, b, path)
	if err != nil {
		return nil, c.rxErr, err
	}

	defer resp.Body.Close()

	u, err := c.rxSerde.StreamFrom(resp.Body)
	if err != nil {
		return nil, c.rxErr, err
	}

	return resp, u, nil
}

func (c *connImpl[RX, TX]) callRaw(
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
