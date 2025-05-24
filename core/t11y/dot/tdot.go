package dot

import (
	"errors"

	"github.com/hkoosha/giraffe/core/t11y"
)

func M[A any](a A, err error) A {
	return t11y.Must(a, err)
}

func OK(err error) {
	t11y.Ensure(err)
}

// =============================================================================

func T[A any](a A, err error) (A, error) {
	if err == nil {
		return a, nil
	}

	return a, E(err)
}

func E(err ...error) error {
	switch {
	case len(err) == 0:
		return t11y.Traced(nil)

	case len(err) == 1 && err[0] == nil:
		return nil

	case len(err) == 1:
		return t11y.Traced(err[0])

	default:
		return t11y.Traced(errors.Join(err...))
	}
}

func EF(format string, v ...any) error {
	return t11y.TracedFmt(format, v...)
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

// =============================================================================

func N(name string, v any) t11y.Named {
	return t11y.Of(name, v)
}

// =============================================================================

func R[T any](t T) *T {
	return &t
}

func D[V any](v *V) V {
	if v == nil {
		var res V
		return res
	}

	return *v
}

func Copy[T any](
	t *T,
) *T {
	if t == nil {
		return t
	}

	return R(D(t))
}
