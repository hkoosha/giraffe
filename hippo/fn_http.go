package hippo

import (
	"maps"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	. "github.com/hkoosha/giraffe/internal/dot1"
)

const (
	httpInputEndpoint = "endpoint"
	httpInputPath     = "path"

	httpInputHeader   = "header"
	httpInputBody     = "body"
	httpInputUrlQuery = "query"
	httpInputMethod   = "method"
	httpInputOkCodes  = "ok_codes"

	httpOutputStatus  = "status"
	httpOutputBody    = "body"
	httpOutputHeaders = "headers"
)

var (
	nameRe = regexp.MustCompile("^[a-zA-Z0-9-_]+$")
	addrRe = regexp.MustCompile(`^(http|https)://(?P<addr>[a-zA-Z0-9-_.]{1,255})(:(?P<port>\d{1,5}))?$`)

	addrReNames = slices.DeleteFunc(addrRe.SubexpNames()[1:], func(it string) bool {
		return it == ""
	})
)

type HttpFn struct {
	cnx       conn.Conn[any]
	endpoints map[string]string
}

func (e *HttpFn) Fn() *Fn {
	return MustFnOf(e.exe).
		WithInput(
			Q(httpInputEndpoint),
			Q(httpInputPath),
		).
		WithOptional(
			Q(httpInputHeader),
			Q(httpInputBody),
			Q(httpInputUrlQuery),
			Q(httpInputMethod),
			Q(httpInputOkCodes),
		).
		WithOutput(
			Q(httpOutputStatus),
			Q(httpOutputBody),
			Q(httpOutputHeaders),
		)
}

func (e *HttpFn) WithConn(
	cnx conn.Conn[any],
) *HttpFn {
	cp := e.shallow()
	cp.cnx = cnx
	return cp
}

func (e *HttpFn) WithEndpoints(
	endpoints map[string]string,
) *HttpFn {
	for name, addr := range endpoints {
		if !nameRe.MatchString(name) {
			panic(EF("invalid endpoint name: %s", name))
		}
		if !addrRe.MatchString(addr) {
			panic(EF("invalid endpoint address: %s", addr))
		}
		matches := addrRe.FindStringSubmatch(addr)
		groups := make(map[string]string)
		for i, n := range addrReNames {
			if i != 0 && n != "" {
				groups[name] = matches[i]
			}
		}

		if strings.Contains(groups["address"], "..") {
			panic(EF("invalid endpoint address: %s", addr))
		}

		if groups["port"] != "" {
			port := M(strconv.Atoi(groups["port"]))
			if port < 1 || 65534 < port {
				panic(EF("invalid endpoint port: %s", addr))
			}
		}
	}

	cp := e.shallow()
	cp.endpoints = maps.Clone(endpoints)
	return cp
}

func (e *HttpFn) shallow() *HttpFn {
	cp := *e
	cp.endpoints = maps.Clone(e.endpoints)
	return &cp
}

func (e *HttpFn) getEndpoint(
	dat giraffe.Datum,
) (string, error) {
	qvEndpointName := M(dat.QStr(httpInputEndpoint))

	endpoint, ok := e.endpoints[qvEndpointName]
	if !ok {
		return "", EF("missing endpoint: %s", qvEndpointName)
	}

	return endpoint, nil
}

func (e *HttpFn) getPath(
	dat giraffe.Datum,
	endpoint string,
) (string, error) {
	pathParts := []string{endpoint}

	for _, part := range strings.Split(M(dat.QStr(httpInputPath)), "/") {
		switch {
		case strings.HasPrefix(part, ":"):
			pValue, err := dat.Query(part[1:])
			if err != nil {
				return "", err
			}
			if pValue.Type().IsInt() {
				pathParts = append(pathParts, M(pValue.Int()).String())
			} else if pValue.Type().IsStr() {
				pathParts = append(pathParts, M(pValue.Str()))
			} else {
				panic("todo")
			}

		default:
			pathParts = append(pathParts, part)
		}
	}

	return conn.Join(pathParts...), nil
}

func (e *HttpFn) getUrlQuery(
	dat giraffe.Datum,
) (string, error) {
	if !dat.Has(Q(httpInputUrlQuery)) {
		return "", nil
	}

	kv, err := dat.QKv(httpInputUrlQuery)
	if err != nil {
		return "", err
	}

	uQueries := make([]string, 0, len(kv))
	for k, v := range kv {
		uQueries = append(uQueries, k+"="+v)
	}

	return strings.Join(uQueries, "&"), nil
}

func (e *HttpFn) getHeaders(
	dat giraffe.Datum,
) (map[string]string, error) {
	if dat.Has(Q(httpInputHeader)) {
		return map[string]string{}, nil
	}

	return dat.QKv(httpInputHeader)
}

func (e *HttpFn) getBody(
	dat giraffe.Datum,
) ([]byte, error) {
	if !dat.Has(Q(httpInputBody)) {
		return nil, nil
	}

	b, err := dat.Query(httpInputBody)
	if err != nil {
		return nil, err
	}

	return b.MarshalJSON()
}

func (e *HttpFn) getMethod(
	dat giraffe.Datum,
	hasBody bool,
) (string, error) {
	if !dat.Has(Q(httpInputMethod)) {
		if hasBody {
			return http.MethodPost, nil
		} else {
			return http.MethodGet, nil
		}
	}

	// TODO prevent get with body?

	return dat.QStr(httpInputMethod)
}

func (e *HttpFn) exe(
	ctx Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	endpoint, err := e.getEndpoint(dat)
	if err != nil {
		return OfErr(), err
	}

	path, err := e.getPath(dat, endpoint)
	if err != nil {
		return OfErr(), err
	}

	uQuery, err := e.getUrlQuery(dat)
	if err != nil {
		return OfErr(), err
	}
	if len(uQuery) > 0 {
		path += "?" + uQuery
	}

	headers, err := e.getHeaders(dat)
	if err != nil {
		return OfErr(), err
	}

	body, err := e.getBody(dat)
	if err != nil {
		return OfErr(), err
	}

	method, err := e.getMethod(dat, len(body) > 0)
	if err != nil {
		return OfErr(), err
	}

	// TODO ok codes

	cnx := e.cnx.
		Cfg().
		WithHeaderOverwrites(true, headers).
		WithMethod(method).
		Conn()

	resp, err := cnx.Call(ctx, body, path)
	if err != nil {
		return OfErr(), err
	}

	return giraffe.FromJsonable(resp)
}
