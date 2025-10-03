package hippo

import (
	"net/http"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/g11y"
)

type httpFn struct {
	cnx    conn.Raw
	method string
}

func (e *httpFn) ekran(
	ctx Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	panic("todo")

	// var resp any
	//
	// if err != nil {
	// 	return OfErr(), err
	// }
	//
	// return giraffe.FromJsonable(resp)
}

func HttpFn(
	name string,
	cnx conn.Raw,
) *Fn {
	g11y.NonNil(cnx)

	fn := httpFn{
		cnx:    cnx,
		method: http.MethodGet,
	}

	return MustFnOf(fn.ekran).
		Named(name).
		Dump()
}
