package dot

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/t11y"
	"github.com/hkoosha/giraffe/core/t11y/dot"
)

// ================================================== ERROR MANAGEMENT SHORTCUTS.

// E trace error(s).
func E(err ...error) error {
	return dot.E(err...)
}

// EF trace formated error.
func EF(format string, v ...any) error {
	return dot.EF(format, v...)
}

// M Must.
func M[A any](a A, err error) A {
	return dot.M(a, err)
}

// N Named.
func N(name string, v any) t11y.Named {
	return dot.N(name, v)
}

func OK(err error) {
	dot.OK(err)
}

// T Traced
func T[A any](a A, err error) (A, error) {
	return dot.T[A](a, err)
}

// ======================================================= DATUM, QUERY, TUPLES.

// Q query.
func Q(
	q string,
) giraffe.Query {
	return giraffe.Q(q)
}

// P pair.
func P[V giraffe.Safe](
	q giraffe.Query,
	v V,
) giraffe.Tuple {
	return giraffe.TupleOf(q, v)
}

func Dat[V giraffe.Safe](v V) giraffe.Datum {
	return giraffe.Of(v)
}

func Dat1[V giraffe.Safe](
	q giraffe.Query,
	v V,
) giraffe.Datum {
	return giraffe.Of1(q, v)
}

func DatN(
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

// ======================================================================== PTR.

// R reference
func R[T any](
	t T,
) *T {
	return dot.R[T](t)
}

// D dereference
func D[T any](
	t *T,
) T {
	return dot.D[T](t)
}

func Clone[T any](
	t *T,
) *T {
	if t == nil {
		return t
	}

	return R(D(t))
}
