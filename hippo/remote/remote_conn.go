package remote

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/hippo"
)

const EkranPath = "/ekran"

func Remote(
	plan string,
	cnx conn.Raw,
) *hippo.Fn {
	cfg := cnx.
		Cfg().
		WithPathPrefix(EkranPath).
		WithSerdes(RequestSerde(), giraffe.DatumSerde())

	fn := remoteFn{
		cnx:  conn.Make[Request, giraffe.Datum](cfg),
		plan: plan,
	}

	return hippo.MustFnCtxOf(fn.Ekran)
}

type remoteFn struct {
	cnx  conn.Conn[Request, giraffe.Datum]
	plan string
}

func (m *remoteFn) String() string {
	return "RemoteFn"
}

func (m *remoteFn) Ekran(
	ctx hippo.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	return m.cnx.Call(ctx, Request{
		Init:          dat,
		Plan:          m.plan,
		Compensations: nil,
	})
}
