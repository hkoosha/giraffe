package hippo

import (
	"maps"
	"net/url"
	"regexp"
	"strings"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/conn/httpmethod"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
)

var httpMethodFnNameRe *regexp.Regexp

func init() {
	all := make([]string, 0)
	for _, m := range httpmethod.All() {
		all = append(all, m.String())
	}

	sb := strings.Builder{}
	sb.WriteString("^(?P<method>")
	sb.WriteString(strings.Join(all, "|"))
	sb.WriteString(`)#\d+$`)

	httpMethodFnNameRe = regexp.MustCompile(sb.String())
}

func validateEndpoint(
	e string,
) error {
	u, err := url.Parse(e)
	if err != nil {
		return E(err)
	}

	ep := u.Scheme + "://" + u.Host
	if p := u.Port(); p != "" {
		ep += ":" + p
	}

	if ep != e {
		return EF("endpoint must be only scheme, host and port: %s", e)
	}

	return nil
}

func (d *DatumTunnels) clone() *DatumTunnels {
	endpoints := make(map[httpmethod.T]map[string]string)
	for k, v := range d.endpoints {
		endpoints[k] = maps.Clone(v)
	}

	return &DatumTunnels{
		cnx:        d.cnx,
		endpoints:  endpoints,
		dynHeaders: maps.Clone(d.dynHeaders),
		headers:    maps.Clone(d.headers),
	}
}

func (d *DatumTunnels) with(
	name string,
	endpoint string,
	method httpmethod.T,
) (*DatumTunnels, error) {
	// TODO validate name

	if err := validateEndpoint(endpoint); err != nil {
		return nil, err
	}

	cp := d.clone()
	cp.endpoints[method][name] = endpoint
	return cp, nil
}

func (d *DatumTunnels) getMethod(
	call Call,
) (httpmethod.T, error) {
	var err error
	mStr := ""
	switch {
	case M(call.Args().Has("method")):
		mStr, err = call.Args().QStr("method")
		if err != nil {
			return httpmethod.GET, err
		}

	case httpMethodFnNameRe.MatchString(call.Name()):
		mStr = httpMethodFnNameRe.FindStringSubmatch(call.Name())[1]
	}

	if mStr == "" {
		return httpmethod.GET, EF("http method not set")
	}
	method, ok := httpmethod.Of(mStr)
	if !ok {
		return httpmethod.GET, EF("unsupported http method: %s", mStr)
	}

	return method, nil
}

func (d *DatumTunnels) getHeaders(
	call Call,
) (map[string]string, error) {
	if ok, err := call.Args().Has("headers"); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	kv, err := call.Args().QKv("headers")
	if err != nil {
		return nil, err
	}

	return kv, nil
}

func (d *DatumTunnels) mkTun(
	name string,
	endpoint string,
	method httpmethod.T,
	h map[string]string,
) *DatumTunnel {
	f := MkTunnel(
		name,
		d.cnx.Cfg().
			WithEndpoint(endpoint).
			WithMethod(method.String()).
			Datum(),
	)

	hs := d.headers
	if hs == nil {
		hs = h
	} else if h != nil {
		for k, v := range h {
			hs[k] = v
		}
	}

	if hs != nil {
		f = f.WithEnforcedHeaders(hs)
	}

	if d.dynHeaders != nil {
		f = f.WithGlobalHeaders(d.dynHeaders)
	}

	if method.HasBody() {
		f = f.WithBody()
	}

	return f
}

func (d *DatumTunnels) exe(
	ctx gtx.Context,
	call Call,
) (giraffe.Datum, error) {
	endpoint, err := call.Args().QStr("channel")
	if err != nil {
		return dErr, err
	}

	path, err := call.Args().QStr("path")
	if err != nil {
		return dErr, err
	}

	_, err = url.Parse(conn.Join("http://localhost", path))
	if err != nil {
		return dErr, E(err)
	}

	method, err := d.getMethod(call)
	if err != nil {
		return dErr, err
	}

	ep, ok := d.endpoints[method][endpoint]
	if !ok {
		return dErr, EF("missing endpoint: %s", endpoint)
	}

	h, err := d.getHeaders(call)
	if err != nil {
		return dErr, err
	}

	tun := d.mkTun(call.Name(), ep, method, h)

	hasExtra, err := call.Args().Has("extra")
	if err != nil {
		return dErr, err
	}

	if hasExtra {
		extra, eErr := call.Args().Get("extra")
		if eErr != nil {
			return dErr, eErr
		}

		tun = tun.WithExtra(extra)
	}

	return tun.exe(ctx, call)
}
