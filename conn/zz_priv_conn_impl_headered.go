package conn

import (
	"context"
	"net/http"
)

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) HPatch(
	ctx context.Context,
	body TX,
	path ...string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	const m = http.MethodPatch
	return c.hCall(ctx, m, &body, path)
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) HPut(
	ctx context.Context,
	body TX,
	path ...string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	const m = http.MethodPut
	return c.hCall(ctx, m, &body, path)
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) HPost(
	ctx context.Context,
	body TX,
	path ...string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	const m = http.MethodPost
	return c.hCall(ctx, m, &body, path)
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) HGet(
	ctx context.Context,
	path ...string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	const m = http.MethodGet
	return c.hCall(ctx, m, nil, path)
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) HDelete(
	ctx context.Context,
	path ...string,
) (
	headers map[string]string,
	_ RX,
	_ error,
) {
	const m = http.MethodDelete
	return c.hCall(ctx, m, nil, path)
}

//nolint:nonamedreturns
func (c *connImpl[TX, RX]) HCall(
	ctx context.Context,
	body *TX,
	path ...string,
) (
	status int,
	headers map[string]string,
	_ RX,
	_ error,
) {
	return c.hsCall(ctx, c.cfg.http.defaultMethod, body, path)
}
