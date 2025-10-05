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

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) Patch(
	ctx context.Context,
	body TX,
	path ...string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	const m = http.MethodPatch
	return c.call(ctx, m, &body, path)
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) Put(
	ctx context.Context,
	body TX,
	path ...string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	const m = http.MethodPut
	return c.call(ctx, m, &body, path)
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) Post(
	ctx context.Context,
	body TX,
	path ...string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	const m = http.MethodPost
	return c.call(ctx, m, &body, path)
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) Get(
	ctx context.Context,
	path ...string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	const m = http.MethodGet
	return c.call(ctx, m, nil, path)
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) Delete(
	ctx context.Context,
	path ...string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	const m = http.MethodDelete
	return c.call(ctx, m, nil, path)
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) Call(
	ctx context.Context,
	body *TX,
	path ...string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	return c.call(ctx, c.cfg.http.defaultMethod, body, path)
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) call(
	ctx context.Context,
	method string,
	reqBody *TX,
	path []string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	var serialized []byte
	var err error

	switch {
	case reqBody != nil:
		serialized, err = c.txSerde.Write(*reqBody)

	default:
		serialized, err = nil, nil
	}

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

	headers = make(map[string]string, len(resp.Header))
	for k, v := range resp.Header {
		if len(v) > 1 {
			panic(EF("todo multivalued headers: %v", resp.Header))
		}
		if len(v) == 0 {
			panic(EF("todo: empty headers: %v", resp.Header))
		}

		headers[k] = v[0]
	}

	return headers, u, nil
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
