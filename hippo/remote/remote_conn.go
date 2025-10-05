package remote

import (
	"encoding/json"
	"io"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/hippo"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
)

const EkranPath = "/ekran"

//goland:noinspection GrazieInspection

func Remote(
	url string,
	plan string,
	cnx conn.Raw,
) *hippo.Fn {
	cfg := cnx.Cfg().
		WithEndpoint(url).
		WithPathPrefix(EkranPath).
		WithSerdes(RequestSerde(), giraffe.DatumSerde())

	fn := remoteFn{
		cnx:  cfg.Conn(),
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

func decode(
	r io.Reader,
) (giraffe.Datum, error) {
	var res any
	dec := json.NewDecoder(r)
	dec.UseNumber()
	dec.DisallowUnknownFields()
	if err := dec.Decode(&res); err != nil {
		return giraffe.OfErr(), hippoerr.NewRemoteError(
			"failed to decode response",
			err,
		)
	}

	dat, err := giraffe.From(res)
	if err != nil {
		return giraffe.OfErr(), hippoerr.NewRemoteError(
			"failed to decode response",
			err,
		)
	}

	return dat, nil
}

func (m *remoteFn) Ekran(
	ctx hippo.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	payload := Request{
		Init:          dat,
		Plan:          m.plan,
		Compensations: nil,
	}

	return m.cnxTo.Call(ctx, payload)

	decoded, dErr := decode(resp.Body)
	if dErr != nil {
		return giraffe.OfErr(), dErr
	}

	return decoded, nil
}
