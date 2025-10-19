package hippo

import (
	"maps"
	"net/http"
	"slices"
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

var bodiless = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodDelete:  {},
	http.MethodConnect: {},
	http.MethodOptions: {},
	http.MethodTrace:   {},
}

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
	return FnOf(e.exe).
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
	if !M(dat.Has(Q(HttpInputUrlQuery))) {
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
	if M(dat.Has(Q(HttpInputHeader))) {
		return map[string]string{}, nil
	}

	return dat.QKv(HttpInputHeader)
}

func (e *HttpFn) getBody(
	dat giraffe.Datum,
) (giraffe.Datum, int, error) {
	if !M(dat.Has(Q(HttpInputBody))) {
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
	if !M(dat.Has(Q(HttpInputMethod))) {
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

	headers, rx, err := cfg.Datum().HCall(ctx, bodyR, path...)
	if err != nil {
		return OfErr(), err
	}

	ret := map[string]any{
		HttpOutputHeaders: headers,
		HttpOutputBody:    rx,
	}

	return giraffe.FromJsonable(ret)
}

// =============================================================================

func MkHttpCallFn(
	cnx conn.Datum,
	endpoint string,
	path string,
) (*HttpCallFn, error) {
	fn := &HttpCallFn{
		cnx:      cnx,
		endpoint: endpoint,
		method:   http.MethodGet,
		okCode:   -1,

		origPath:         "fail",
		path:             "fail",
		pathPartsStatic:  nil,
		pathPartsVar:     nil,
		query:            "fail",
		queryPartsStatic: nil,
		queryPartsVar:    nil,

		headers: make(map[string]giraffe.Query, 0),
		body:    nil,
	}

	return fn.WithPath(path)
}

type HttpCallFn struct {
	cnx              conn.Datum
	body             *giraffe.Query
	headers          map[string]giraffe.Query
	query            string
	method           string
	endpoint         string
	origPath         string
	path             string
	pathPartsStatic  []string
	queryPartsStatic []string
	queryPartsVar    []giraffe.Query
	pathPartsVar     []giraffe.Query
	okCode           int
}

func (e *HttpCallFn) WithOkCodes(
	v int,
) *HttpCallFn {
	if e.okCode == v {
		return e
	}

	cp := e.clone()
	cp.okCode = v
	return cp
}

func (e *HttpCallFn) WithoutOkCodes() *HttpCallFn {
	return e.WithOkCodes(-1)
}

func (e *HttpCallFn) WithBody(
	q giraffe.Query,
) *HttpCallFn {
	if e.body != nil && *e.body == q {
		return e
	}

	cp := e.clone()
	cp.body = &q
	return cp
}

func (e *HttpCallFn) WithoutBody() *HttpCallFn {
	if e.body == nil {
		return e
	}

	cp := e.clone()
	cp.body = nil
	return cp
}

func (e *HttpCallFn) WithMethod(
	v string,
) (*HttpCallFn, error) {
	if e.method == v {
		return e, nil
	}

	if _, ok := bodiless[v]; ok && e.body != nil {
		return nil, EF(
			"method cannot have body and a body is set, method=%s body=%s",
			v,
			e.body.String(),
		)
	}

	cp := e.clone()
	cp.method = v
	return cp, nil
}

func (e *HttpCallFn) WithEndpoint(
	v string,
) (*HttpCallFn, error) {
	// TODO validate endpoints

	if e.endpoint == v {
		return e, nil
	}

	cp := e.clone()
	cp.endpoint = v
	return cp, nil
}

func (e *HttpCallFn) WithHeaders(
	v map[string]giraffe.Query,
) *HttpCallFn {
	if maps.Equal(e.headers, v) {
		return e
	}

	cp := e.clone()

	if len(v) == 0 {
		cp.headers = map[string]giraffe.Query{}
	} else {
		cp.headers = maps.Clone(v)
	}

	return cp
}

func (e *HttpCallFn) WithPath(
	path string,
) (*HttpCallFn, error) {
	if e.origPath == path {
		return e, nil
	}

	// TODO use url.Parse or something.
	cp := e.clone()

	pathOrig := path

	path, query, _ := strings.Cut(path, "?")

	pathParts := strings.Split(path, "/")
	pathPartsStatic := make([]string, len(pathParts))
	pathPartsVar := make([]giraffe.Query, len(pathParts))
	pathPartAnyVar := false
	for i, part := range pathParts {
		switch {
		case part == "":
			return nil, EF("empty path part: %s", path)

		case strings.HasPrefix(part, "."):
			q, err := giraffe.GQParse(part)
			if err != nil {
				return nil, err
			}
			pathPartsVar[i] = q
			pathPartAnyVar = true

		default:
			pathPartsStatic[i] = part
		}
	}
	if !pathPartAnyVar {
		pathPartsStatic = nil
		pathPartsVar = nil
	}

	queryParts := strings.Split(query, "&")
	queryPartsStatic := make([]string, len(queryParts))
	queryPartsVar := make([]giraffe.Query, len(queryParts))
	queryPartsAnyVar := false
	for i, part := range queryParts {
		switch {
		case part == "":
			return nil, EF("empty query part: %s", path)

		case strings.HasPrefix(part, "."):
			q, err := giraffe.GQParse(part)
			if err != nil {
				return nil, err
			}
			queryPartsVar[i] = q
			queryPartsAnyVar = true

		default:
			queryPartsStatic[i] = part
		}
	}
	if !queryPartsAnyVar {
		queryPartsStatic = nil
		queryPartsVar = nil
	}

	cp.pathPartsStatic = pathPartsStatic
	cp.pathPartsVar = pathPartsVar
	cp.queryPartsStatic = queryPartsStatic
	cp.queryPartsVar = queryPartsVar
	cp.path = path
	cp.query = query
	cp.origPath = pathOrig

	return cp, nil
}

func (e *HttpCallFn) Fn() *Fn {
	fn := FnOf(e.exe).
		WithOutput(
			Q(HttpOutputBody),
			Q(HttpOutputHeaders),
		)

	for _, d := range e.headers {
		fn = fn.WithInput(d)
	}

	if e.body != nil {
		fn = fn.WithInput(*e.body)
	}

	for _, d := range e.pathPartsVar {
		fn = fn.WithInput(d)
	}

	return fn
}

func (e *HttpCallFn) clone() *HttpCallFn {
	return &HttpCallFn{
		cnx:              e.cnx,
		method:           e.method,
		pathPartsStatic:  slices.Clone(e.pathPartsStatic),
		pathPartsVar:     slices.Clone(e.pathPartsVar),
		queryPartsStatic: slices.Clone(e.queryPartsStatic),
		queryPartsVar:    slices.Clone(e.queryPartsVar),
		path:             e.path,
		query:            e.query,
		origPath:         e.origPath,
		endpoint:         e.endpoint,
		body:             e.body,
		headers:          maps.Clone(e.headers),
		okCode:           e.okCode,
	}
}

//nolint:nestif
func (e *HttpCallFn) mkPath(
	dat giraffe.Datum,
) ([]string, error) {
	if e.pathPartsStatic == nil && e.queryPartsStatic == nil {
		return []string{e.origPath}, nil
	}

	parts := make([]string, 0, len(e.pathPartsStatic)+len(e.queryPartsStatic))

	if e.pathPartsStatic == nil {
		for i := range len(e.pathPartsStatic) {
			switch s := e.pathPartsStatic[i]; {
			case s != "":
				parts = append(parts, s)

			default:
				ss, err := dat.Query(e.pathPartsVar[i].String())
				if err != nil {
					return nil, E(
						err,
						EF("missing path variable in context: %s", e.pathPartsVar[i].String()),
					)
				}

				sss, err := ss.SimpleString()
				if err != nil {
					return nil, E(
						err,
						EF("path variable cannot be formatted to string: %s",
							e.pathPartsVar[i].String()),
					)
				}

				parts = append(parts, sss)
			}
		}
	} else {
		parts = append(parts, e.path)
	}

	if e.query != "" {
		parts = append(parts, "?")

		if e.queryPartsStatic == nil {
			for i := range len(e.queryPartsStatic) {
				switch s := e.queryPartsStatic[i]; {
				case s != "":
					parts = append(parts, s)

				default:
					ss, err := dat.Query(e.queryPartsVar[i].String())
					if err != nil {
						return nil, E(
							err,
							EF("missing query variable in context: %s",
								e.queryPartsVar[i].String()),
						)
					}

					sss, err := ss.SimpleString()
					if err != nil {
						return nil, E(
							err,
							EF("query variable cannot be formatted to string: %s",
								e.queryPartsVar[i].String()),
						)
					}

					parts = append(parts, sss)
				}
			}
		} else {
			parts = append(parts, e.query)
		}
	}

	return parts, nil
}

func (e *HttpCallFn) getBody(
	dat giraffe.Datum,
) (giraffe.Datum, int, error) {
	if e.body == nil {
		return OfErr(), 0, nil
	}

	q, err := dat.Query(e.body.String())
	if err != nil {
		return OfErr(), 0, err
	}

	l, err := q.Len()
	if err != nil {
		return OfErr(), 0, err
	}

	return q, l, nil
}

func (e *HttpCallFn) exe(
	ctx Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	path, err := e.mkPath(dat)
	if err != nil {
		return OfErr(), err
	}

	headers := make(map[string]string, len(e.headers))
	for header, q := range e.headers {
		dd, err0 := dat.Query(q.String())
		if err0 != nil {
			return OfErr(), err0
		}

		ddStr, err0 := dd.SimpleString()
		if err0 != nil {
			return OfErr(), err0
		}

		headers[header] = ddStr
	}

	body, l, err := e.getBody(dat)
	if err != nil {
		return OfErr(), err
	}

	cfg := e.cnx.Cfg().
		WithHeaderOverwrites(true, headers).
		WithMethod(e.method)

	if e.okCode > 0 {
		cfg.WithExpectingStatusCode(e.okCode)
	}

	cfg, err = cfg.WithEndpointNamed(e.endpoint)
	if err != nil {
		return OfErr(), err
	}

	bodyR := &body
	if l == 0 {
		bodyR = nil
	}

	headers, rx, err := cfg.Datum().HCall(ctx, bodyR, path...)
	if err != nil {
		return OfErr(), err
	}

	ret := map[string]any{
		HttpOutputHeaders: headers,
		HttpOutputBody:    rx,
	}

	return giraffe.FromJsonable(ret)
}
