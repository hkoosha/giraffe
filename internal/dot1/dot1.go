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

func Of0[V giraffe.SafeType1](v V) giraffe.Datum {
	return giraffe.Of(v)
}

func OfN(
	pairs ...giraffe.Tuple,
) (giraffe.Datum, error) {
	return giraffe.OfN(pairs...)
}

func Q(
	q string,
) giraffe.Query {
	return giraffe.Q(q)
}

func P[V giraffe.SafeType1](
	q giraffe.Query,
	v V,
) giraffe.Tuple {
	return giraffe.TupleOf(q, v)
}
