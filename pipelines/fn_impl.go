package pipelines

import (
	"context"

	"github.com/hkoosha/giraffe"
)

func Static(
	dat giraffe.Datum,
) Fn {
	fn := staticFn{
		dat: dat,
	}

	return fn.Ekran
}

type staticFn struct {
	dat giraffe.Datum
}

func (m staticFn) String() string {
	return "StaticFn[]"
}

func (m staticFn) Ekran(
	_ context.Context,
	_ giraffe.Datum,
) (giraffe.Datum, error) {
	return m.dat, nil
}

// =============================================================================.

func Scoped(
	scope giraffe.Query,
	fn Fn,
) Fn {
	sFn := scopedFn{
		scope: scope,
		fn:    fn,
	}

	return sFn.Ekran
}

type scopedFn struct {
	fn    Fn
	scope giraffe.Query
}

func (m *scopedFn) String() string {
	return "ScopedFn[]"
}

func (m *scopedFn) Ekran(
	ctx context.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	ekran, err := m.fn(ctx, dat)
	if err != nil {
		return giraffe.OfErr(), err
	}

	return giraffe.Of1(m.scope, ekran), nil
}

// =============================================================================.

type FlattenScopeError struct {
	msg string
}

func (e *FlattenScopeError) Error() string {
	return e.msg
}

func FlattenScoped(
	scope giraffe.Query,
	fn Fn,
) Fn {
	sFn := flattenScopedFn{
		scope: scope,
		fn:    fn,
	}

	return sFn.Ekran
}

type flattenScopedFn struct {
	fn    Fn
	scope giraffe.Query
}

func (m *flattenScopedFn) String() string {
	return "FlattenScopedFn[]"
}

func (m *flattenScopedFn) Ekran(
	ctx context.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	ekran, err := m.fn(ctx, dat)
	if err != nil {
		return giraffe.OfErr(), err
	}

	if !ekran.Type().IsObj() {
		return giraffe.OfErr(), &FlattenScopeError{
			msg: "expecting an object to flatten, but got: " + ekran.String(),
		}
	}

	ret := giraffe.OfEmpty()

	it, err := ekran.Iter2()
	if err != nil {
		return giraffe.OfErr(), err
	}

	for k, v := range it {
		// TODO will doubly scape.
		key := giraffe.EscapedQ(k)
		if ret, err = ret.Set(m.scope.Plus(key), v); err != nil {
			return giraffe.OfErr(), err
		}
	}

	return giraffe.Of1(m.scope, ekran), nil
}
