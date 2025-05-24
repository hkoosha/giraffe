package hippo

import (
	"maps"
	"net/url"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/conn/httpmethod"
	"github.com/hkoosha/giraffe/core/t11y"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

func MkTunnels(
	cnx conn.Datum,
) *DatumTunnels {
	t11y.NonNil(cnx)

	return &DatumTunnels{
		cnx: cnx,
		endpoints: map[httpmethod.T]map[string]string{
			httpmethod.GET:    {},
			httpmethod.POST:   {},
			httpmethod.PUT:    {},
			httpmethod.PATCH:  {},
			httpmethod.DELETE: {},
		},
		headers:    nil,
		dynHeaders: nil,
	}
}

type DatumTunnels struct {
	cnx        conn.Datum
	endpoints  map[httpmethod.T]map[string]string
	headers    map[string]string
	dynHeaders map[giraffe.Query]string
}

// =============================================================================

func (d *DatumTunnels) MustWithPost(
	name string,
	endpoint string,
) *DatumTunnels {
	m := httpmethod.POST
	return M(d.with(name, endpoint, m))
}

func (d *DatumTunnels) WithPost(
	name string,
	endpoint string,
) (*DatumTunnels, error) {
	m := httpmethod.POST
	return d.with(name, endpoint, m)
}

func (d *DatumTunnels) WithPostU(
	name string,
	endpoint *url.URL,
) (*DatumTunnels, error) {
	m := httpmethod.POST
	return d.with(name, endpoint.String(), m)
}

func (d *DatumTunnels) MustWithPostU(
	name string,
	endpoint *url.URL,
) *DatumTunnels {
	m := httpmethod.POST
	return M(d.with(name, endpoint.String(), m))
}

// =============================================================================

func (d *DatumTunnels) MustWithPatch(
	name string,
	endpoint string,
) *DatumTunnels {
	m := httpmethod.PATCH
	return M(d.with(name, endpoint, m))
}

func (d *DatumTunnels) WithPatch(
	name string,
	endpoint string,
) (*DatumTunnels, error) {
	m := httpmethod.PATCH
	return d.with(name, endpoint, m)
}

func (d *DatumTunnels) WithPatchU(
	name string,
	endpoint *url.URL,
) (*DatumTunnels, error) {
	m := httpmethod.PATCH
	return d.with(name, endpoint.String(), m)
}

func (d *DatumTunnels) MustWithPatchU(
	name string,
	endpoint *url.URL,
) *DatumTunnels {
	m := httpmethod.PATCH
	return M(d.with(name, endpoint.String(), m))
}

// =============================================================================

func (d *DatumTunnels) MustWithDelete(
	name string,
	endpoint string,
) *DatumTunnels {
	m := httpmethod.DELETE
	return M(d.with(name, endpoint, m))
}

func (d *DatumTunnels) WithDelete(
	name string,
	endpoint string,
) (*DatumTunnels, error) {
	m := httpmethod.DELETE
	return d.with(name, endpoint, m)
}

func (d *DatumTunnels) WithDeleteU(
	name string,
	endpoint *url.URL,
) (*DatumTunnels, error) {
	m := httpmethod.DELETE
	return d.with(name, endpoint.String(), m)
}

func (d *DatumTunnels) MustWithDeleteU(
	name string,
	endpoint *url.URL,
) *DatumTunnels {
	m := httpmethod.DELETE
	return M(d.with(name, endpoint.String(), m))
}

// =============================================================================

func (d *DatumTunnels) MustWithPut(
	name string,
	endpoint string,
) *DatumTunnels {
	m := httpmethod.PUT
	return M(d.with(name, endpoint, m))
}

func (d *DatumTunnels) WithPut(
	name string,
	endpoint string,
) (*DatumTunnels, error) {
	m := httpmethod.PUT
	return d.with(name, endpoint, m)
}

func (d *DatumTunnels) WithPutU(
	name string,
	endpoint *url.URL,
) (*DatumTunnels, error) {
	m := httpmethod.PUT
	return d.with(name, endpoint.String(), m)
}

func (d *DatumTunnels) MustWithPutU(
	name string,
	endpoint *url.URL,
) *DatumTunnels {
	m := httpmethod.PUT
	return M(d.with(name, endpoint.String(), m))
}

// =============================================================================

func (d *DatumTunnels) MustWithGet(
	name string,
	endpoint string,
) *DatumTunnels {
	m := httpmethod.GET
	return M(d.with(name, endpoint, m))
}

func (d *DatumTunnels) WithGet(
	name string,
	endpoint string,
) (*DatumTunnels, error) {
	m := httpmethod.GET
	return d.with(name, endpoint, m)
}

func (d *DatumTunnels) WithGetU(
	name string,
	endpoint *url.URL,
) (*DatumTunnels, error) {
	m := httpmethod.GET
	return d.with(name, endpoint.String(), m)
}

func (d *DatumTunnels) MustWithGetU(
	name string,
	endpoint *url.URL,
) *DatumTunnels {
	m := httpmethod.GET
	return M(d.with(name, endpoint.String(), m))
}

// =============================================================================

func (d *DatumTunnels) WithEnforcedHeaders(
	h map[string]string,
) *DatumTunnels {
	cp := d.clone()
	cp.headers = maps.Clone(h)
	return cp
}

func (d *DatumTunnels) WithDynHeaders(
	h map[giraffe.Query]string,
) *DatumTunnels {
	cp := d.clone()
	cp.dynHeaders = maps.Clone(h)
	return cp
}

func (d *DatumTunnels) WithoutDynHeaders() *DatumTunnels {
	return d.WithDynHeaders(map[giraffe.Query]string{})
}

// =============================================================================

func (d *DatumTunnels) RegisterTo(
	reg *FnRegistry,
) (*FnRegistry, error) {
	var err error

	for method := range d.endpoints {
		reg, err = reg.WithNamed(method.String(), d.Fn())
		if err != nil {
			return nil, err
		}
	}

	return reg, nil
}

func (d *DatumTunnels) Fn() *Fn {
	return FnOf(d.exe).
		WithOutput(HttpOutputBody).
		WithOutput(HttpOutputHeaders)
}
