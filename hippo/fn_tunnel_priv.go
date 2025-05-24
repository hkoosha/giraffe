package hippo

import (
	"maps"
	"slices"
	"strings"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn/headers"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
)

func simple(
	dat giraffe.Datum,
	query giraffe.Query,
) (string, error) {
	ss, err := dat.Get(query)
	if err != nil {
		return "", E(
			err,
			EF("missing path variable in context: %s", query.String()),
		)
	}

	sss, err := ss.SimpleString()
	if err != nil {
		return "", E(
			err,
			EF("path variable cannot be formatted to string: %s", query.String()),
		)
	}

	return sss, nil
}

type datumTunnelPath struct {
	fullPath         string
	pathOnly         string
	queryOnly        string
	pathPartsStatic  []string
	pathPartsVar     []giraffe.Query
	queryPartsStatic []string
	queryPartsVar    []giraffe.Query
}

func (d *datumTunnelPath) clone() *datumTunnelPath {
	return &datumTunnelPath{
		fullPath:         d.fullPath,
		pathOnly:         d.pathOnly,
		queryOnly:        d.queryOnly,
		pathPartsStatic:  slices.Clone(d.pathPartsStatic),
		pathPartsVar:     slices.Clone(d.pathPartsVar),
		queryPartsStatic: slices.Clone(d.queryPartsStatic),
		queryPartsVar:    slices.Clone(d.queryPartsVar),
	}
}

func (d *datumTunnelPath) withPath(
	path string,
) (*datumTunnelPath, error) {
	cp := d.clone()
	cp.fullPath = path

	if err := adjustPath(cp); err != nil {
		return nil, err
	}

	if err := adjustQuery(cp); err != nil {
		return nil, err
	}

	cp, err := cp.optimized()
	if err != nil {
		return nil, err
	}

	return cp, nil
}

func (d *datumTunnelPath) optimized() (*datumTunnelPath, error) {
	if d.pathPartsStatic == nil && d.queryPartsStatic == nil {
		//nolint:exhaustruct
		return &datumTunnelPath{
			fullPath: d.fullPath,
		}, nil
	} else {
		return d, nil
	}
}

func adjustPath(
	d *datumTunnelPath,
) error {
	if d.fullPath == "" {
		panic("full path not set")
	}

	// TODO use url.Parse or something.
	pathOnly, _, _ := strings.Cut(d.fullPath, "?")
	pathOnly, _ = strings.CutPrefix(pathOnly, "/")

	pathParts := strings.Split(pathOnly, "/")
	pathPartsStatic := make([]string, len(pathParts))
	pathPartsVar := make([]giraffe.Query, len(pathParts))
	pathPartAnyVar := false
	for i, part := range pathParts {
		switch {
		case part == "":
			return EF("empty path part: %s", d.fullPath)

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
		pathOnly = ""
	}

	d.pathPartsStatic = pathPartsStatic
	d.pathPartsVar = pathPartsVar
	d.pathOnly = pathOnly

	return nil
}

func adjustQuery(
	d *datumTunnelPath,
) error {
	if d.fullPath == "" {
		panic("full path not set")
	}

	// TODO use url.Parse or something.
	_, queryOnly, _ := strings.Cut(d.fullPath, "?")

	queryParts := strings.Split(queryOnly, "&")
	queryPartsStatic := make([]string, len(queryParts))
	queryPartsVar := make([]giraffe.Query, len(queryParts))
	queryPartsAnyVar := false

	if queryOnly != "" {
		for i, part := range queryParts {
			switch {
			case part == "":
				return EF("empty query part: %s", d.fullPath)

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
		queryOnly = ""
	}

	d.queryPartsStatic = queryPartsStatic
	d.queryPartsVar = queryPartsVar
	d.queryOnly = queryOnly

	return nil
}

func (d *datumTunnelPath) mkPath(
	dat giraffe.Datum,
) ([]string, error) {
	if d.pathOnly == "" && d.queryOnly == "" {
		return []string{d.fullPath}, nil
	}

	parts := make(
		[]string,
		0,
		len(d.pathPartsStatic)+len(d.queryPartsStatic),
	)

	if d.pathPartsVar != nil {
		for i := range len(d.pathPartsStatic) {
			switch s := d.pathPartsStatic[i]; {
			case s != "":
				parts = append(parts, s)

			default:
				ss, err := simple(dat, d.pathPartsVar[i])
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
	if d.queryPartsVar != nil {
		parts = append(parts, "?")

		if d.queryPartsStatic != nil {
			for i := range len(d.queryPartsStatic) {
				switch s := d.queryPartsStatic[i]; {
				case s != "":
					parts = append(parts, s)

				default:
					ss, err := simple(dat, d.queryPartsVar[i])
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

func (h *DatumTunnel) clone() *DatumTunnel {
	return &DatumTunnel{
		cnx:             h.cnx,
		enforcedHeaders: maps.Clone(h.enforcedHeaders),
		hasBody:         h.hasBody,
		name:            h.name,
		template:        h.template.clone(),
		globalHeaders:   maps.Clone(h.globalHeaders),
	}
}

func (h *DatumTunnel) mkHeaders(
	call Call,
) (map[string]string, error) {
	hh := maps.Clone(h.enforcedHeaders)

	ok, err := call.Args().Has("bearer")
	if err != nil {
		return nil, err
	} else if ok {
		tokenQStr, tErr := call.Args().QStr("bearer")
		if tErr != nil {
			return nil, tErr
		}

		tokenQ, tErr := giraffe.GQParse(tokenQStr)
		if tErr != nil {
			return nil, tErr
		}

		token, tErr := call.Data().QStr(tokenQ)
		if tErr != nil {
			return nil, tErr
		}

		hh[headers.Authorization] = "Bearer " + token
	}

	ok, err = call.Data().Has(HttpInputHeader)
	if err != nil {
		return nil, err
	}

	if ok {
		kv, err := call.Data().QKv(HttpInputHeader)
		if err != nil {
			return nil, err
		}

		for k, v := range kv {
			hh[k] = v
		}
	}

	return hh, nil
}

func (h *DatumTunnel) mkPath(
	dat giraffe.Datum,
) ([]string, error) {
	pathTpl, err := dat.QStr(giraffe.Q(HttpInputPath))
	if err != nil {
		return nil, err
	}

	tpl, err := h.template.withPath(pathTpl)
	if err != nil {
		return nil, err
	}

	path, err := tpl.mkPath(dat)
	if err != nil {
		return nil, err
	}

	return path, nil
}

func (h *DatumTunnel) getBody(
	call Call,
) (*giraffe.Datum, error) {
	if !h.hasBody {
		return nil, nil
	}

	if ok, err := call.Args().Has("no_body"); err != nil {
		return nil, err
	} else if ok {
		if ok, err = call.Args().QBln("no_body"); err != nil {
			return nil, err
		} else if ok {
			return nil, nil
		}
	}

	body, err := call.Data().Get(giraffe.Q(HttpInputBody))
	if err != nil {
		return nil, err
	}

	if h.extra != nil {
		body, err = body.Merge(*h.extra)
		if err != nil {
			return nil, err
		}
	}

	l, err := body.Len()
	if err != nil {
		return nil, err
	}

	if l == 0 {
		return nil, nil
	}

	return &body, nil
}

func (h *DatumTunnel) exe(
	ctx gtx.Context,
	call Call,
) (giraffe.Datum, error) {
	path, err := h.mkPath(call.Args())
	if err != nil {
		return dErr, err
	}

	reqHeaders, err := h.mkHeaders(call)
	if err != nil {
		return dErr, err
	}

	body, err := h.getBody(call)
	if err != nil {
		return dErr, err
	}

	cnx := h.cnx.Cfg().AndHeaders(reqHeaders).Datum()

	status, respHeaders, rx, err := cnx.HCall(ctx, body, path...)
	if err != nil {
		return dErr, err
	}

	isErr := !cnx.IsExpected(ctx, status)
	ret := map[string]any{
		HttpOutputHeaders: respHeaders,
		HttpOutputBody:    rx,
		HttpOutputStatus:  status,
		HttpOutputError:   isErr,
	}

	if isErr {
		diag := rx.MarshalJsonString()
		return dErr, EF(
			"unexpected status: code=%d\nheaders=%v\nbody=%s",
			status,
			respHeaders,
			diag,
		)
	}

	return giraffe.FromJsonable(ret)
}
