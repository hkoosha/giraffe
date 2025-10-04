package dot

import (
	"errors"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/hippo"
	"github.com/hkoosha/giraffe/t11y"
	"github.com/hkoosha/giraffe/zebra/zptr"
)

// ================================================== ERROR MANAGEMENT SHORTCUTS.

// T traced error if error, result if otherwise
func T[A any](a A, err error) (A, error) {
	if err != nil {
		return a, E(err)
	}

	return a, nil
}

// E trace error(s).
func E(err ...error) error {
	switch {
	case len(err) == 0:
		return t11y.Traced(nil)

	case len(err) == 1:
		return t11y.Traced(err[0])

	default:
		return t11y.Traced(errors.Join(err...))
	}
}

// EF trace formated error.
func EF(format string, v ...any) error {
	return t11y.TracedFmt(format, v...)
}

// M Must.
func M[A any](
	a A,
	err error,
) A {
	return t11y.Must(a, err)
}

// N Named.
func N(
	name string,
	v any,
) t11y.Named {
	return t11y.Of(name, v)
}

func OK(
	err error,
) {
	t11y.Ensure(err)
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

func Of0[V giraffe.Safe](v V) giraffe.Datum {
	return giraffe.Of(v)
}

func Of1[V giraffe.Safe](
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

// ========================================================================= FN.

type Fn = *hippo.Fn

func OfFn(
	exe hippo.Exe,
) Fn {
	return hippo.MustFnOf(exe)
}

// ======================================================================== PTR.

func Ref[T any](
	t T,
) *T {
	return zptr.Ref[T](t)
}

func Deref[T any](
	t *T,
) T {
	return zptr.Deref[T](t)
}
