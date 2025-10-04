package dot1

import (
	"github.com/hkoosha/giraffe"
)

func OfErr() giraffe.Datum {
	return giraffe.OfErr()
}

func OfEmpty() giraffe.Datum {
	return giraffe.OfEmpty()
}

func Of0[V giraffe.Safe](v V) giraffe.Datum {
	return giraffe.Of(v)
}

func OfN(
	pairs ...giraffe.Tuple,
) (giraffe.Datum, error) {
	return giraffe.OfN(pairs...)
}

func Q(
	q string,
) giraffe.GQuery {
	return giraffe.Q(q)
}

func P[V giraffe.Safe](
	q giraffe.GQuery,
	v V,
) giraffe.Tuple {
	return giraffe.TupleOf(q, v)
}
