package hippo

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/g11y/gtx"
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/typing"
)

// ============================================================================.

type Exe0 = func(giraffe.Datum) (giraffe.Datum, error)

type ExeCtx = func(
	gtx.Context,
	giraffe.Datum,
) (giraffe.Datum, error)

type Exe = func(
	gtx.Context,
	giraffe.Datum,
) (giraffe.Datum, error)

// ============================================================================.

func Static(
	dat giraffe.Datum,
) *Fn_ {
	return M(FnOf(func(
		gtx.Context,
		giraffe.Datum,
	) (giraffe.Datum, error) {
		return dat, nil
	}))
}

// ============================================================================.

func MustFnOf(
	exe Exe,
) *Fn_ {
	return M(FnOf(exe))
}

func FnOf(
	exe Exe,
) (*Fn_, error) {
	t := typing.OfVirtual()

	fn := &Fn_{
		exe:           exe,
		scopedOut:     "",
		scopedIn:      "",
		inputs:        nil,
		outputs:       nil,
		optionals:     nil,
		replicated:    nil,
		selected:      nil,
		swapped:       nil,
		noOverwriting: false,
		typ:           t,
		name:          "#" + t.String(),
	}

	var err error = nil
	if !fn.IsValid() {
		err = E(errInvalidFn)
	}

	//nolint:nilnil
	return fn, err
}

func MustFnCtxOf(
	exe ExeCtx,
) *Fn_ {
	return M(FnCtxOf(exe))
}

func FnCtxOf(
	exeCtx ExeCtx,
) (*Fn_, error) {
	exe := func(
		ctx gtx.Context,
		dat giraffe.Datum,
	) (giraffe.Datum, error) {
		return exeCtx(ctx, dat)
	}

	return FnOf(exe)
}

func MustFnOf0(
	exe0 Exe0,
) *Fn_ {
	return M(FnOf0(exe0))
}

func FnOf0(
	exe0 Exe0,
) (*Fn_, error) {
	exe := func(
		_ gtx.Context,
		dat giraffe.Datum,
	) (giraffe.Datum, error) {
		return exe0(dat)
	}

	return FnOf(exe)
}
