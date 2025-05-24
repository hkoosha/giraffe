package conn

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/hkoosha/giraffe/conn/internal"
	"github.com/hkoosha/giraffe/core/serdes"
	"github.com/hkoosha/giraffe/core/t11y"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/internal/reflected"
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

func (c *connImpl[TX, RX]) call(
	ctx context.Context,
	method string,
	reqBody *TX,
	path []string,
) (RX, error) {
	_, rx, err := c.hCall(ctx, method, reqBody, path)
	return rx, err
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) hCall(
	ctx context.Context,
	method string,
	reqBody *TX,
	path []string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	_, headers, rx, err := c.hsCall(ctx, method, reqBody, path)
	return headers, rx, err
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) hsCall(
	ctx context.Context,
	method string,
	reqBody *TX,
	path []string,
) (
	status int,
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
		return 0, nil, c.rxErr, E(err)
	}

	resp, err := c.callRaw(ctx, method, bytes.NewReader(serialized), path)
	if err != nil {
		return 0, nil, c.rxErr, err
	}

	//goland:noinspection GoMaybeNil - false positive
	defer resp.Body.Close()

	u, err := c.rxSerde.StreamFrom(resp.Body)
	if err != nil {
		return 0, nil, c.rxErr, err
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

	return resp.StatusCode, headers, u, nil
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

	for k, v := range c.cfg.header.overwrite {
		req.Header.Set(k, v)
	}

	for k, v := range c.cfg.header.overwriters {
		req.Header.Set(k, v(ctx, c.cfg))
	}

	resp, err := c.std.Do(req)
	if err != nil {
		return nil, E(err)
	}

	return resp, nil
}
