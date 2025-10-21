package hippo

import (
	"maps"
	"strings"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/conn/httpmethod"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

func MkDatumChannel(
	cnx conn.Datum,
	name string,
	path string,
) (DatumChannel, error) {
	if cnx.Cfg().Endpoint() == "" {
		panic(EF("http endpoint not set on channel connection"))
	}

	// TODO validate name as simple machine name.

	dc := DatumChannel{
		cnx:     cnx,
		headers: make(map[string]struct{}),
		hasBody: false,
		name:    name,
		path:    "",
		dcPath: datumChannelPath{
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

type datumChannelPath struct {
	pathOnly         string
	queryOnly        string
	pathPartsStatic  []string
	pathPartsVar     []giraffe.Query
	queryPartsStatic []string
	queryPartsVar    []giraffe.Query
}

func (d *datumChannelPath) isZero() bool {
	return d.pathOnly == ""
}

func (d *datumChannelPath) hasDynPath() bool {
	return d.pathPartsVar != nil
}

func (d *datumChannelPath) hasDynQuery() bool {
	return d.queryPartsVar != nil
}

func (d *datumChannelPath) mkPath(
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

type DatumChannel struct {
	cnx     conn.Datum
	headers map[string]struct{}
	name    string
	path    string
	dcPath  datumChannelPath
	hasBody bool
}

func (h *DatumChannel) Id() string {
	return h.name
}

// =====================================

func (h *DatumChannel) WithHeaders(
	name string,
	headers map[string]struct{},
) DatumChannel {
	cp := h.shallow(name)
	cp.headers = headers
	return cp
}

func (h *DatumChannel) WithBody(
	name string,
) DatumChannel {
	return h.SetHasBody(name, true)
}

func (h *DatumChannel) WithoutBody(
	name string,
) DatumChannel {
	return h.SetHasBody(name, false)
}

func (h *DatumChannel) SetHasBody(
	name string,
	b bool,
) DatumChannel {
	if b && !httpmethod.MustOf(h.cnx.Cfg().Method()).HasBody() {
		panic(EF("http method does not take body: %s", h.cnx.Cfg().Method()))
	}

	cp := h.shallow(name)
	cp.hasBody = b
	return cp
}

func (h *DatumChannel) WithPath(
	name string,
	path string,
) (DatumChannel, error) {
	cp := h.shallow(name)
	cp.path = path

	// TODO use url.Parse or something.
	pathOnly, queryOnly, _ := strings.Cut(path, "?")

	if err := procPath(&cp, pathOnly); err != nil {
		return datumChannelErr, err
	}
	if err := procQuery(&cp, queryOnly); err != nil {
		return datumChannelErr, err
	}

	if cp.dcPath.pathPartsStatic == nil && cp.dcPath.queryPartsStatic == nil {
		cp.dcPath = datumChannelPathZero
	}

	return cp, nil
}

// =====================================

func procPath(
	dc *DatumChannel,
	pathOnly string,
) error {
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
	dc *DatumChannel,
	queryOnly string,
) error {
	queryParts := strings.Split(queryOnly, "&")
	queryPartsStatic := make([]string, len(queryParts))
	queryPartsVar := make([]giraffe.Query, len(queryParts))
	queryPartsAnyVar := false
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
	if !queryPartsAnyVar {
		queryPartsStatic = nil
		queryPartsVar = nil
	}

	dc.dcPath.queryPartsStatic = queryPartsStatic
	dc.dcPath.queryPartsVar = queryPartsVar
	dc.dcPath.queryOnly = queryOnly

	return nil
}

func (h *DatumChannel) shallow(
	name string,
) DatumChannel {
	return DatumChannel{
		cnx:     h.cnx,
		headers: maps.Clone(h.headers),
		hasBody: h.hasBody,
		path:    h.path,
		name:    name,
		dcPath:  h.dcPath,
	}
}

func (h *DatumChannel) mkHeaders(
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
var datumChannelErr = DatumChannel{}

//nolint:exhaustruct
var datumChannelPathZero = datumChannelPath{}
