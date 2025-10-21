package hippo

import (
	"maps"
	"strings"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/conn/httpmethod"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

func MustMkTunnel(
	name string,
	cnx conn.Datum,
	path string,
) *DatumTunnel {
	return M(MkTunnel(name, cnx, path))
}

func MkTunnel(
	name string,
	cnx conn.Datum,
	path string,
) (*DatumTunnel, error) {
	if cnx.Cfg().Endpoint() == "" {
		panic(EF("http endpoint not set on tunnel connection"))
	}

	// TODO validate name as simple machine name.

	dc := &DatumTunnel{
		cnx:     cnx,
		headers: make(map[string]struct{}),
		hasBody: false,
		name:    name,
		path:    "",
		dcPath: datumTunnelPath{
			pathPartsStatic:  nil,
			pathPartsVar:     nil,
			queryPartsStatic: nil,
			queryPartsVar:    nil,
			pathOnly:         "",
			queryOnly:        "",
		},
	}

	return dc.WithPath(dc.name, path)
}

// =====================================

type simpleStr = func(giraffe.Query) (string, error)

type datumTunnelPath struct {
	pathOnly         string
	queryOnly        string
	pathPartsStatic  []string
	pathPartsVar     []giraffe.Query
	queryPartsStatic []string
	queryPartsVar    []giraffe.Query
}

func (d *datumTunnelPath) isZero() bool {
	return d.pathOnly == ""
}

func (d *datumTunnelPath) hasDynPath() bool {
	return d.pathPartsVar != nil
}

func (d *datumTunnelPath) hasDynQuery() bool {
	return d.queryPartsVar != nil
}

func (d *datumTunnelPath) mkPath(
	get simpleStr,
) ([]string, error) {
	parts := make(
		[]string,
		0,
		len(d.pathPartsStatic)+len(d.queryPartsStatic),
	)

	if d.hasDynPath() {
		for i := range len(d.pathPartsStatic) {
			switch s := d.pathPartsStatic[i]; {
			case s != "":
				parts = append(parts, s)

			default:
				q := d.pathPartsVar[i]
				ss, err := get(q)
				if err != nil {
					return nil, err
				}

				parts = append(parts, ss)
			}
		}
	} else {
		parts = append(parts, d.pathOnly)
	}

	//nolint:nestif
	if d.hasDynQuery() {
		parts = append(parts, "?")

		if d.queryPartsStatic != nil {
			for i := range len(d.queryPartsStatic) {
				switch s := d.queryPartsStatic[i]; {
				case s != "":
					parts = append(parts, s)

				default:
					q := d.queryPartsVar[i]
					ss, err := get(q)
					if err != nil {
						return nil, err
					}

					parts = append(parts, ss)
				}
			}
		} else {
			parts = append(parts, d.queryOnly)
		}
	}

	return parts, nil
}

// =====================================

type DatumTunnel struct {
	cnx     conn.Datum
	headers map[string]struct{}
	name    string
	path    string
	dcPath  datumTunnelPath
	hasBody bool
}

func (h *DatumTunnel) Id() string {
	return h.name
}

func (h *DatumTunnel) Fn() *Fn {
	return mkHttpCallFn(h).Fn()
}

// =====================================

func (h *DatumTunnel) WithHeaders(
	name string,
	headers map[string]struct{},
) *DatumTunnel {
	cp := h.shallow(name)
	cp.headers = headers
	return cp
}

func (h *DatumTunnel) WithBody(
	name string,
) *DatumTunnel {
	return h.SetHasBody(name, true)
}

func (h *DatumTunnel) WithoutBody(
	name string,
) *DatumTunnel {
	return h.SetHasBody(name, false)
}

func (h *DatumTunnel) SetHasBody(
	name string,
	b bool,
) *DatumTunnel {
	if b && !httpmethod.MustOf(h.cnx.Cfg().Method()).HasBody() {
		panic(EF("http method does not take body: %s", h.cnx.Cfg().Method()))
	}

	cp := h.shallow(name)
	cp.hasBody = b
	return cp
}

func (h *DatumTunnel) WithPath(
	name string,
	path string,
) (*DatumTunnel, error) {
	cp := h.shallow(name)
	cp.path = path

	// TODO use url.Parse or something.
	pathOnly, queryOnly, _ := strings.Cut(path, "?")

	if err := procPath(cp, pathOnly); err != nil {
		return nil, err
	}
	if err := procQuery(cp, queryOnly); err != nil {
		return nil, err
	}

	if cp.dcPath.pathPartsStatic == nil && cp.dcPath.queryPartsStatic == nil {
		cp.dcPath = datumTunnelPathZero
	}

	return cp, nil
}

// =====================================

func procPath(
	dc *DatumTunnel,
	pathOnly string,
) error {
	pathOnly, _ = strings.CutPrefix(pathOnly, "/")

	pathParts := strings.Split(pathOnly, "/")
	pathPartsStatic := make([]string, len(pathParts))
	pathPartsVar := make([]giraffe.Query, len(pathParts))
	pathPartAnyVar := false
	for i, part := range pathParts {
		switch {
		case part == "":
			return EF("empty path part: %s", dc.path)

		case strings.HasPrefix(part, "."):
			q, err := giraffe.GQParse(part)
			if err != nil {
				return err
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

	dc.dcPath.pathPartsStatic = pathPartsStatic
	dc.dcPath.pathPartsVar = pathPartsVar
	dc.dcPath.pathOnly = pathOnly

	return nil
}

func procQuery(
	dc *DatumTunnel,
	queryOnly string,
) error {
	queryParts := strings.Split(queryOnly, "&")
	queryPartsStatic := make([]string, len(queryParts))
	queryPartsVar := make([]giraffe.Query, len(queryParts))
	queryPartsAnyVar := false

	if queryOnly != "" {
		for i, part := range queryParts {
			switch {
			case part == "":
				return EF("empty query part: %s", dc.path)

			case strings.HasPrefix(part, "."):
				q, err := giraffe.GQParse(part)
				if err != nil {
					return err
				}
				queryPartsVar[i] = q
				queryPartsAnyVar = true

			default:
				queryPartsStatic[i] = part
			}
		}
	}
	if !queryPartsAnyVar {
		queryPartsStatic = nil
		queryPartsVar = nil
	}

	dc.dcPath.queryPartsStatic = queryPartsStatic
	dc.dcPath.queryPartsVar = queryPartsVar
	dc.dcPath.queryOnly = queryOnly

	return nil
}

func (h *DatumTunnel) shallow(
	name string,
) *DatumTunnel {
	return &DatumTunnel{
		cnx:     h.cnx,
		headers: maps.Clone(h.headers),
		hasBody: h.hasBody,
		path:    h.path,
		name:    name,
		dcPath:  h.dcPath,
	}
}

func (h *DatumTunnel) mkHeaders(
	get simpleStr,
) (map[string]string, error) {
	headers := make(map[string]string, len(h.headers))

	for header := range h.headers {
		q, err := giraffe.GQParse(header)
		if err != nil {
			return nil, err
		}

		d, err := get(q)
		if err != nil {
			return nil, err
		}

		headers[header] = d
	}

	return headers, nil
}

//nolint:exhaustruct
var datumTunnelPathZero = datumTunnelPath{}
