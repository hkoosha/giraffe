package ghttp

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/zebra/serdes"
)

type noBodyT struct{}

func (noBodyT) Read([]byte) (int, error)         { return 0, io.EOF }
func (noBodyT) Close() error                     { return nil }
func (noBodyT) WriteTo(io.Writer) (int64, error) { return 0, nil }

var nobody = noBodyT{}

// ============================================================================.

func newConn[T, U any](
	cfg *config,
	tSerde serdes.Serde[T],
	uSerde serdes.Serde[U],
) *conn[T, U] {
	g11y.NonNil(cfg, tSerde, uSerde)
	cfg.Ensure()

	return &conn[T, U]{
		cfg:    cfg,
		std:    cfg.Std(),
		tSerde: tSerde,
		uSerde: uSerde,
	}
}

type conn[T any, U any] struct {
	cfg    *config
	std    *http.Client
	tSerde serdes.Serde[T]
	uSerde serdes.Serde[U]
	tErr   T
	uErr   U
}

func (c *conn[T, U]) Std() *http.Client {
	return c.cfg.Std()
}

func (c *conn[T, U]) Cfg() Config {
	return c.cfg
}

func (c *conn[T, U]) mapErrResp(
	resp *http.Response,
) (*http.Response, error) {
	if c.cfg.resp.expectStatusCode == 0 ||
		resp.StatusCode == c.cfg.resp.expectStatusCode {
		return resp, nil
	}

	var err error
	var body []byte
	if resp.Body != nil {
		_ = resp.Body.Close()
	}

	_ = err
	_ = body

	_ = &FailedResponseError{
		Reason: ReasonUnexpectedStatusCode,
	}

	panic("todo")
}

func (c *conn[T, U]) Patch(
	ctx context.Context,
	body T,
	path ...string,
) (U, error) {
	const m = http.MethodPatch
	return c.call(ctx, m, body, path)
}

func (c *conn[T, U]) Put(
	ctx context.Context,
	body T,
	path ...string,
) (U, error) {
	const m = http.MethodPut
	return c.call(ctx, m, body, path)
}

func (c *conn[T, U]) Post(
	ctx context.Context,
	body T,
	path ...string,
) (U, error) {
	const m = http.MethodPost
	return c.call(ctx, m, body, path)
}

func (c *conn[T, U]) Get(
	ctx context.Context,
	path ...string,
) (U, error) {
	const m = http.MethodGet
	return c.call(ctx, m, nobody, path)
}

func (c *conn[T, U]) Delete(
	ctx context.Context,
	path ...string,
) (U, error) {
	const m = http.MethodDelete
	return c.call(ctx, m, nil, path)
}

func (c *conn[T, U]) Call(
	ctx context.Context,
	method string,
	body T,
	path ...string,
) (U, error) {
	return c.call(ctx, method, body, path)
}

func (c *conn[T, U]) call(
	ctx context.Context,
	method string,
	body any,
	path []string,
) (U, error) {
	var b io.Reader

	if body == nil || body == http.NoBody || body == nobody {
		b = nobody
	} else if cast, ok := body.(io.Reader); ok {
		b = cast
	} else if cast, ok := body.([]byte); ok {
		b = bytes.NewReader(cast)
	} else if cast, ok := body.(T); ok {
		serialized, err := c.tSerde.Write(cast)
		if err != nil {
			return c.uErr, err
		}
		b = bytes.NewReader(serialized)
	} else {
		panic("unreachable, unknown body type: " + reflect.TypeOf(body).String())
	}

	resp, err := c.callRaw(ctx, method, b, path)
	if err != nil {
		return c.uErr, err
	}

	if resp.Body == nil {
		return c.uSerde.Read([]byte{})
	}

	defer resp.Body.Close()

	u, err := c.uSerde.ReadFrom(resp.Body)
	if err != nil {
		return c.uErr, err
	}

	return u, nil
}

func (c *conn[T, U]) callRaw(
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
