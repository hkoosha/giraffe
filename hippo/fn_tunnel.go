package hippo

import (
	"maps"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/cmd"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/conn/httpmethod"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

const (
	HttpInputPath   = "path"
	HttpInputHeader = "headers"
	HttpInputBody   = "body"

	HttpOutputBody    = "body"
	HttpOutputHeaders = "headers"
)

func MkTunnel(
	name string,
	cnx conn.Datum,
) *DatumTunnel {
	if cnx.Cfg().Endpoint() == "" {
		panic(EF("http endpoint not set on tunnel connection"))
	}

	// TODO validate name as simple machine name.

	dc := &DatumTunnel{
		cnx:             cnx,
		enforcedHeaders: make(map[string]string),
		globalHeaders:   make(map[giraffe.Query]string),
		hasBody:         false,
		name:            name,
		template: &datumTunnelPath{
			pathPartsStatic:  nil,
			pathPartsVar:     nil,
			queryPartsStatic: nil,
			queryPartsVar:    nil,
			pathOnly:         "",
			queryOnly:        "",
			fullPath:         "",
		},
	}

	return dc
}

type DatumTunnel struct {
	cnx             conn.Datum
	name            string
	hasBody         bool
	enforcedHeaders map[string]string
	globalHeaders   map[giraffe.Query]string
	template        *datumTunnelPath
}

func (h *DatumTunnel) Id() string {
	return h.name
}

func (h *DatumTunnel) Fn() *Fn {
	fn := FnOf(h.exe)

	for d := range h.globalHeaders {
		fn = fn.WithInput(d)
	}

	if h.hasBody {
		fn = fn.WithInput(giraffe.Q(HttpInputBody))
	}

	for _, d := range h.template.pathPartsVar {
		fn = fn.WithInput(giraffe.Q("var" + cmd.Sep.String() + d.String()))
	}
	for _, d := range h.template.queryPartsVar {
		fn = fn.WithInput(giraffe.Q("var" + cmd.Sep.String() + d.String()))
	}

	return fn.
		WithOptional(HttpInputHeader).
		WithInput(HttpInputPath).
		WithOutput(HttpOutputBody).
		WithOutput(HttpOutputHeaders)
}

func (h *DatumTunnel) WithEnforcedHeaders(
	headers map[string]string,
) *DatumTunnel {
	cp := h.clone()
	cp.enforcedHeaders = maps.Clone(headers)
	return cp
}

func (h *DatumTunnel) WithoutEnforcedHeaders() *DatumTunnel {
	return h.WithEnforcedHeaders(map[string]string{})
}

func (h *DatumTunnel) WithGlobalHeaders(
	headers map[giraffe.Query]string,
) *DatumTunnel {
	cp := h.clone()
	cp.globalHeaders = maps.Clone(headers)
	return cp
}

func (h *DatumTunnel) WithoutGlobalHeaders() *DatumTunnel {
	return h.WithGlobalHeaders(map[giraffe.Query]string{})
}

func (h *DatumTunnel) WithBody() *DatumTunnel {
	return h.SetHasBody(true)
}

func (h *DatumTunnel) WithoutBody() *DatumTunnel {
	return h.SetHasBody(false)
}

func (h *DatumTunnel) SetHasBody(
	b bool,
) *DatumTunnel {
	if b && !httpmethod.MustOf(h.cnx.Cfg().Method()).HasBody() {
		panic(EF("http method does not take body: %s", h.cnx.Cfg().Method()))
	}

	cp := h.clone()
	cp.hasBody = b
	return cp
}
