package conn

import (
	"context"
	"net/http"
	"slices"
)

func (c *connImpl[TX, RX]) Patch(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, error) {
	const m = http.MethodPatch
	return c.call(ctx, m, &body, path)
}

func (c *connImpl[TX, RX]) Put(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, error) {
	const m = http.MethodPut
	return c.call(ctx, m, &body, path)
}

func (c *connImpl[TX, RX]) Post(
	ctx context.Context,
	body TX,
	path ...string,
) (RX, error) {
	const m = http.MethodPost
	return c.call(ctx, m, &body, path)
}

func (c *connImpl[TX, RX]) Get(
	ctx context.Context,
	path ...string,
) (RX, error) {
	const m = http.MethodGet
	return c.call(ctx, m, nil, path)
}

func (c *connImpl[TX, RX]) Delete(
	ctx context.Context,
	path ...string,
) (RX, error) {
	const m = http.MethodDelete
	return c.call(ctx, m, nil, path)
}

func (c *connImpl[TX, RX]) Call(
	ctx context.Context,
	body *TX,
	path ...string,
) (RX, error) {
	return c.call(ctx, c.cfg.http.defaultMethod, body, path)
}

func (c *connImpl[TX, RX]) IsExpected(
	_ context.Context,
	code int,
) bool {
	return c.cfg.resp.expectStatusCode == nil ||
		slices.Contains(c.cfg.resp.expectStatusCode, code)
}
