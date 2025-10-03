package hippo

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/g11y"
	. "github.com/hkoosha/giraffe/internal/dot1"
)

type httpFn struct {
	method string
	cnx    conn.Raw
}

func (e *httpFn) ekran(
	ctx Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	panic("todo")

	var err error
	var resp any

	if err != nil {
		return OfErr(), err
	}

	return giraffe.FromJsonable(resp)
}

func HttpFn(
	name string,
	cnx conn.Raw,
) *Fn {
	g11y.NonNil(cnx)

	fn := httpFn{
		cnx: cnx,
	}

	return MustFnOf(fn.ekran).
		Named(name).
		Dump()
}
