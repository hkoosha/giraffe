package dot1

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/t11y"
	"github.com/hkoosha/giraffe/core/t11y/dot"
)

func OfErr() giraffe.Datum {
	return giraffe.OfErr()
}

func OfEmpty() giraffe.Datum {
	return giraffe.OfEmpty()
}

func Of[V giraffe.Safe](v V) giraffe.Datum {
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

func P[V giraffe.Safe](
	q giraffe.Query,
	v V,
) giraffe.Tuple {
	return giraffe.TupleOf(q, v)
}

// =============================================================================

func M[A any](a A, err error) A {
	return dot.M(a, err)
}

func OK(err error) {
	dot.OK(err)
}

func E(err ...error) error {
	return dot.E(err...)
}

func EF(format string, v ...any) error {
	return dot.EF(format, v...)
}

func Assert(condition bool) {
	dot.Assert(condition)
}

//goland:noinspection SpellCheckingInspection
func Assertf(condition bool, format string, v ...any) {
	dot.Assertf(condition, format, v...)
}

func N(name string, v any) t11y.Named {
	return dot.N(name, v)
}
