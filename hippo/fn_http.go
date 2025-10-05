package hippo

import (
	"net/http"
	"strings"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	. "github.com/hkoosha/giraffe/internal/dot1"
)

const (
	HttpInputEndpoint = "endpoint"
	HttpInputPath     = "path"

	HttpInputHeader   = "headers"
	HttpInputBody     = "body"
	HttpInputUrlQuery = "query"
	HttpInputMethod   = "method"
	HttpInputOkCodes  = "ok_codes"

	HttpOutputBody    = "body"
	HttpOutputHeaders = "headers"
)

func MkHttpFn(
	cnx conn.Datum,
) *HttpFn {
	return &HttpFn{
		cnx: cnx.Cfg(),
	}
}

type HttpFn struct {
	cnx conn.Config
}

func (e *HttpFn) WithConn(
	cnx conn.Datum,
) *HttpFn {
	cp := e.shallow()
	cp.cnx = cnx.Cfg()
	return cp
}

func (e *HttpFn) Fn() *Fn {
	return MustFnOf(e.exe).
		WithInput(
			Q(HttpInputEndpoint),
			Q(HttpInputPath),
		).
		WithOptional(
			Q(HttpInputHeader),
			Q(HttpInputBody),
			Q(HttpInputUrlQuery),
			Q(HttpInputMethod),
			Q(HttpInputOkCodes),
		).
		WithOutput(
			Q(HttpOutputBody),
			Q(HttpOutputHeaders),
		)
}

// =============================================================================

func (e *HttpFn) shallow() *HttpFn {
	cp := *e
	return &cp
}

func (e *HttpFn) getEndpoint(
	dat giraffe.Datum,
) (string, error) {
	return dat.QStr(HttpInputEndpoint)
}

func (e *HttpFn) getPath(
	dat giraffe.Datum,
) ([]string, error) {
	var parts []string

	for _, part := range strings.Split(M(dat.QStr(HttpInputPath)), "/") {
		switch {
		case strings.HasPrefix(part, ":"):
			pValue, err := dat.Query(part[1:])
			if err != nil {
				return nil, err
			}
			if pValue.Type().IsInt() {
				parts = append(parts, M(pValue.Int()).String())
			} else if pValue.Type().IsStr() {
				parts = append(parts, M(pValue.Str()))
			} else {
				panic("todo")
			}

		default:
			parts = append(parts, part)
		}
	}

	return parts, nil
}

func (e *HttpFn) getUrlQuery(
	dat giraffe.Datum,
) ([]string, error) {
	if !dat.Has(Q(HttpInputUrlQuery)) {
		return nil, nil
	}

	kv, err := dat.QKv(HttpInputUrlQuery)
	if err != nil {
		return nil, err
	}

	uQueries := make([]string, 0, len(kv))
	for k, v := range kv {
		uQueries = append(uQueries, k+"="+v)
	}

	return uQueries, nil
}

func (e *HttpFn) getHeaders(
	dat giraffe.Datum,
) (map[string]string, error) {
	if dat.Has(Q(HttpInputHeader)) {
		return map[string]string{}, nil
	}

	return dat.QKv(HttpInputHeader)
}

func (e *HttpFn) getBody(
	dat giraffe.Datum,
) (giraffe.Datum, int, error) {
	if !dat.Has(Q(HttpInputBody)) {
		return OfErr(), 0, nil
	}

	q, err := dat.Query(HttpInputBody)
	if err != nil {
		return OfErr(), 0, err
	}

	l, err := q.Len()
	if err != nil {
		return OfErr(), 0, err
	}

	return q, l, nil
}

func (e *HttpFn) getMethod(
	dat giraffe.Datum,
	hasBody bool,
) (string, error) {
	if !dat.Has(Q(HttpInputMethod)) {
		if hasBody {
			return http.MethodPost, nil
		} else {
			return http.MethodGet, nil
		}
	}

	// TODO prevent get with body?

	return dat.QStr(HttpInputMethod)
}

func (e *HttpFn) exe(
	ctx Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	path, err := e.getPath(dat)
	if err != nil {
		return OfErr(), err
	}

	if uQuery, uqErr := e.getUrlQuery(dat); uqErr != nil {
		return OfErr(), uqErr
	} else if len(uQuery) > 0 {
		path = append(path, "?")
		path = append(path, uQuery...)
	}

	headers, err := e.getHeaders(dat)
	if err != nil {
		return OfErr(), err
	}

	body, l, err := e.getBody(dat)
	if err != nil {
		return OfErr(), err
	}

	method, err := e.getMethod(dat, l > 0)
	if err != nil {
		return OfErr(), err
	}

	endpoint, err := e.getEndpoint(dat)
	if err != nil {
		return OfErr(), err
	}

	// TODO ok codes

	cfg, err := e.cnx.
		WithHeaderOverwrites(true, headers).
		WithMethod(method).
		WithEndpointNamed(endpoint)
	if err != nil {
		return OfErr(), err
	}

	bodyR := &body
	if l == 0 {
		bodyR = nil
	}

	headers, rx, err := cfg.Datum().Call(ctx, bodyR, path...)
	if err != nil {
		return OfErr(), err
	}

	ret := map[string]any{
		HttpOutputHeaders: headers,
		HttpOutputBody:    rx,
	}

	return giraffe.FromJsonable(ret)
}
