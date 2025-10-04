package dot0

import (
	"errors"
	"slices"

	"github.com/hkoosha/giraffe/t11y"
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

func N(name string, v any) t11y.Named {
	return t11y.Of(name, v)
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

func OK(
	err error,
) {
	t11y.Ensure(err)
}
