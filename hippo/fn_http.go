package hippo

import (
	"maps"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	. "github.com/hkoosha/giraffe/internal/dot0"
	. "github.com/hkoosha/giraffe/internal/dot1"
)

var (
	httpPathVarRe    = regexp.MustCompile("^(:)?([a-zA-Z0-9-_]+)$")
	httpSimpleNameRe = regexp.MustCompile("^[a-zA-Z0-9-_]+$")
)

type HttpFn struct {
	cnx        conn.Conn[any]
	urlQueries map[string]giraffe.Query
	urlParts   []string
}

func (e *HttpFn) shallow() *HttpFn {
	cp := *e
	cp.urlQueries = maps.Clone(e.urlQueries)
	cp.urlParts = slices.Clone(e.urlParts)
	return &cp
}

func (e *HttpFn) WithConn(
	cnx conn.Conn[any],
) *HttpFn {
	cp := e.shallow()
	cp.cnx = cnx
	return cp
}

func (e *HttpFn) WithPath(
	path ...string,
) *HttpFn {
	var flat []string

	for _, p := range path {
		for _, s := range strings.Split(p, "/") {
			if !httpPathVarRe.MatchString(s) {
				panic(EF("invalid path: %v", path))
			}
			flat = append(flat, s)
		}
	}

	cp := e.shallow()
	cp.urlParts = flat
	return cp
}

func (e *HttpFn) WithUrlQueries(
	queries map[string]giraffe.Query,
) *HttpFn {
	for name := range queries {
		if !httpSimpleNameRe.MatchString(name) {
			panic(EF("invalid http queries: %s, %v", name, queries))
		}
	}

	cp := e.shallow()
	cp.urlQueries = maps.Clone(queries)
	return cp
}

func (e *HttpFn) exe(
	ctx Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	path := make([]string, 0, len(e.urlParts)+1+len(e.urlQueries))
	for _, v := range e.urlParts {
		switch {
		case v[0] == ':':
			vq, err := dat.Query(v[1:])
			if err != nil {
				return OfErr(), err
			}
			vqStr, err := vq.FmtStr()
			if err != nil {
				return OfErr(), err
			}
			path = append(path, vqStr)

		default:
			path = append(path, v)
		}
	}

	path = append(path, "?")
	for name, q := range e.urlQueries {
		str, err := dat.QFmtStr(q)
		if err != nil {
			return OfErr(), err
		}
		path = append(path, name+"="+str)
	}

	body, err := e.cnx.Call(ctx, http.NoBody, path...)
	if err != nil {
		return OfErr(), err
	}

	return giraffe.FromJsonable(body)
}

func (e *HttpFn) Fn() *Fn {
	return MustFnOf(e.exe)
}
