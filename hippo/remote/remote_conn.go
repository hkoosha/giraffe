package remote

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	"github.com/hkoosha/giraffe/hippo"
)

const EkranPath = "/ekran"

func Remote(
	plan string,
	cnx conn.Raw,
) *hippo.Fn {
	cfg := cnx.Cfg().WithPathPrefix(EkranPath)

	fn := remoteFn{
		cnx: conn.Make[Request, giraffe.Datum](
			cfg,
			RequestSerde(),
			giraffe.DatumSerde(),
		),
		plan: plan,
	}

	return hippo.FnOf(fn.Ekran)
}

type remoteFn struct {
	cnx  conn.Conn[Request, giraffe.Datum]
	plan string
}

func (m *remoteFn) String() string {
	return "RemoteFn"
}

func (m *remoteFn) Ekran(
	ctx gtx.Context,
	call hippo.Call,
) (giraffe.Datum, error) {
	_, rx, err := m.cnx.HCall(ctx, &Request{
		Init:          call.Data(),
		Plan:          m.plan,
		Compensations: nil,
	})
	return rx, err
}
