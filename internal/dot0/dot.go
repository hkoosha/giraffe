package dot0

import (
	"errors"
	"maps"
	"slices"

	"github.com/hkoosha/giraffe/g11y/named"
	"github.com/hkoosha/giraffe/g11y/t11y"
	"github.com/hkoosha/giraffe/zebra/z"
)

func M[A any](a A, err error) A {
	return t11y.Must(a, err)
}

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

func EF(format string, v ...any) error {
	return t11y.TracedFmt(format, v...)
}

func N(name string, v any) named.Named {
	return named.Of(name, v)
}

func Assert(condition bool) {
	Assertf(condition, "unexpected state")
}

//goland:noinspection SpellCheckingInspection
func Assertf(
	condition bool,
	format string,
	v ...any,
) {
	if !condition {
		panic(EF(format, v...))
	}
}

func Ref[T any](t T) *T {
	return &t
}

func Appended[S ~[]E, E any](s S, e ...E) S {
	return append(slices.Clone(s), e...)
}

func TryAnd[M ~map[K]V, K comparable, V any](
	m M,
	k K,
	v V,
) (M, bool) {
	if _, ok := m[k]; ok {
		return nil, false
	}

	mCp := maps.Clone(m)
	mCp[k] = v

	return mCp, true
}

func Apply[U, V any](
	it []U,
	fn func(U) V,
) []V {
	return z.Applied(it, fn)
}

func OK(
	err error,
) {
	t11y.Ensure(err)
}
