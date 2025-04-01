package dot

import (
	"errors"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/g11y"
)

// ================================================== ERROR MANAGEMENT SHORTCUTS.

// E trace error(s).
func E(err ...error) error {
	switch {
	case len(err) == 0:
		return g11y.Traced(nil)

	case len(err) == 1:
		return g11y.Traced(err[0])

	default:
		return g11y.Traced(errors.Join(err...))
	}
}

// EF trace formated error.
func EF(format string, v ...any) error {
	return g11y.TracedFmt(format, v...)
}

// M Must.
func M[A any](
	a A,
	err error,
) A {
	return g11y.Must(a, err)
}

func OK(
	err error,
) {
	g11y.Ensure(err)
}

// ======================================================= DATUM, QUERY, TUPLES.

// Q query.
func Q(
	q string,
) giraffe.Query {
	return giraffe.Q(q)
}

// P pair.
func P[V giraffe.SafeType1](
	q giraffe.Query,
	v V,
) giraffe.Tuple {
	return giraffe.TupleOf(q, v)
}

func Of0[V giraffe.SafeType1](v V) giraffe.Datum {
	return giraffe.Of(v)
}

func Of1[V giraffe.SafeType1](
	q giraffe.Query,
	v V,
) giraffe.Datum {
	return giraffe.Of1(q, v)
}

func OfN(
	pairs ...giraffe.Tuple,
) (giraffe.Datum, error) {
	return giraffe.OfN(pairs...)
}

func OfErr() giraffe.Datum {
	return giraffe.OfErr()
}

func OfEmpty() giraffe.Datum {
	return giraffe.OfEmpty()
}
