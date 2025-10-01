package hippo

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
	"github.com/hkoosha/giraffe/hippo/internal/privnames"
	. "github.com/hkoosha/giraffe/internal/dot0"
	. "github.com/hkoosha/giraffe/internal/dot1"
	"github.com/hkoosha/giraffe/typing"
)

var errInvalidFn = errors.New("invalid fn")

type Exe0 = func(giraffe.Datum) (giraffe.Datum, error)

type ExeCtx = func(
	context.Context,
	giraffe.Datum,
) (giraffe.Datum, error)

type Exe = func(
	Context,
	giraffe.Datum,
) (giraffe.Datum, error)

// =============================================================================.

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
		exe:        exe,
		scoped:     "",
		inputs:     nil,
		outputs:    nil,
		optionals:  nil,
		replicated: nil,
		selected:   nil,
		swapped:    nil,
		typ:        t,
		name:       "#" + t.String(),
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
		ctx Context,
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
		_ Context,
		dat giraffe.Datum,
	) (giraffe.Datum, error) {
		return exe0(dat)
	}

	return FnOf(exe)
}

type Fn_ struct {
	exe        Exe
	replicated map[giraffe.Query]giraffe.Query
	swapped    map[giraffe.Query]giraffe.Query
	scoped     giraffe.Query
	name       string
	inputs     []giraffe.Query
	outputs    []giraffe.Query
	optionals  []giraffe.Query
	selected   []giraffe.Query
	typ        typing.Type
}

func (f *Fn_) ensure() {
	if !f.IsValid() {
		panic(EF("invalid fn"))
	}
}

func (f *Fn_) Type() typing.Type {
	if f == nil {
		return typing.OfErr()
	}

	return f.typ
}

func (f *Fn_) IsValid() bool {
	return f != nil && f.exe != nil && f.typ.IsValid()
}

func (f *Fn_) AndReplicate(
	replicated map[giraffe.Query]giraffe.Query,
) *Fn_ {
	f.ensure()

	replicated = maps.Clone(replicated)
	maps.Copy(replicated, f.replicated)

	clone := f.clone()
	clone.replicated = replicated
	return clone
}

func (f *Fn_) WithReplicated(
	replicated map[giraffe.Query]giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.replicated = maps.Clone(replicated)
	return clone
}

func (f *Fn_) AndSwapping(
	swapping map[giraffe.Query]giraffe.Query,
) *Fn_ {
	f.ensure()

	swapping = maps.Clone(swapping)
	maps.Copy(swapping, f.swapped)

	clone := f.clone()
	clone.swapped = swapping
	return clone
}

func (f *Fn_) WithSwapping(
	swapping map[giraffe.Query]giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.swapped = maps.Clone(swapping)
	return clone
}

func (f *Fn_) AndScope(
	scope giraffe.Query,
) *Fn_ {
	f.ensure()

	return f.WithScope(f.scoped.Plus(scope))
}

func (f *Fn_) WithScope(
	scope giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.scoped = scope
	return clone
}

func (f *Fn_) AndInputs(
	inputs ...giraffe.Query,
) *Fn_ {
	return f.WithInput(append(inputs, f.inputs...)...)
}

func (f *Fn_) WithInput(
	inputs ...giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.inputs = inputs
	return clone
}

func (f *Fn_) AndOptionals(
	optionals ...giraffe.Query,
) *Fn_ {
	return f.WithOptional(append(optionals, f.optionals...)...)
}

func (f *Fn_) WithOptional(
	optionals ...giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.optionals = optionals
	return clone
}

func (f *Fn_) AndOutputs(
	outputs ...giraffe.Query,
) *Fn_ {
	return f.WithOutput(append(outputs, f.outputs...)...)
}

func (f *Fn_) WithOutput(
	outputs ...giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.outputs = slices.Clone(outputs)
	return clone
}

func (f *Fn_) AndSelect(
	select_ ...giraffe.Query,
) *Fn_ {
	return f.Select(append(select_, f.selected...)...)
}

func (f *Fn_) Select(
	select_ ...giraffe.Query,
) *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.selected = slices.Clone(select_)
	return clone
}

func (f *Fn_) SelectAll() *Fn_ {
	f.ensure()

	clone := f.clone()
	clone.selected = nil
	return clone
}

func (f *Fn_) Named(
	name string,
) *Fn_ {
	f.ensure()

	if !privnames.SimpleName.MatchString(name) {
		panic(EF("invalid fn name: %s", name))
	}

	clone := f.clone()
	clone.name = name
	return clone
}

func (f *Fn_) Dump() *Fn_ {
	return f
}

func (f *Fn_) String() string {
	return fmt.Sprintf("Fn[%s][%s]", f.typ, f.name)
}

// =============================================================================.

func chkDatPresent(
	dat giraffe.Datum,
	keys []giraffe.Query,
) error {
	if len(keys) == 0 {
		return nil
	}

	var missing []giraffe.Query
	for _, k := range keys {
		if !dat.Has(k) {
			missing = append(missing, k)
		}
	}

	if len(missing) > 0 {
		return hippoerr.NewFnMissingKeysError(missing)
	}

	return nil
}

// =====================================.

func (f *Fn_) replicate(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	if len(f.replicated) == 0 {
		return dat, nil
	}

	for from, into := range f.replicated {
		if !dat.Has(from) {
			continue
		}

		val, err := dat.Get(from)
		if err != nil {
			return OfErr(), err
		}

		dat, err = dat.Set(into, val)
		if err != nil {
			return OfErr(), err
		}
	}

	return dat, nil
}

func (f *Fn_) scope(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	if f.scoped == "" {
		return dat, nil
	}

	return giraffe.Of1(f.scoped, dat), nil
}

func (f *Fn_) select_(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	if len(f.selected) == 0 {
		return dat, nil
	}

	selected := make(map[giraffe.Query]giraffe.Datum, len(f.selected))
	for _, k := range f.selected {
		if !dat.Has(k) {
			continue
		}
		v, err := dat.Get(k)
		if err != nil {
			return OfErr(), err
		}

		selected[k] = v
	}

	return Of0(selected), nil
}

func (f *Fn_) swap(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	if len(f.swapped) == 0 {
		return dat, nil
	}

	if !dat.Type().IsObj() {
		return OfErr(), EF("expecting an object, got: %s", dat.Type())
	}

	iter, err := dat.Iter2()
	if err != nil {
		return OfErr(), err
	}

	ret := make(map[giraffe.Query]giraffe.Datum, M(dat.Len()))
	for k, v := range iter {
		k := Q(k)
		if swapTo, ok := f.replicated[k]; ok {
			k = swapTo
		}

		ret[k] = v
	}

	return Of0(ret), nil
}

// =====================================.

func (f *Fn_) clone() *Fn_ {
	f.ensure()

	if f == nil {
		return nil
	}

	return &Fn_{
		exe:        f.exe,
		scoped:     f.scoped,
		inputs:     slices.Clone(f.inputs),
		optionals:  slices.Clone(f.optionals),
		outputs:    slices.Clone(f.outputs),
		replicated: maps.Clone(f.replicated),
		swapped:    maps.Clone(f.swapped),
		selected:   slices.Clone(f.selected),
		typ:        f.typ.Clone(),
		name:       f.name,
	}
}

func (f *Fn_) call(
	ctx Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	g11y.NonNil(f, f.exe)

	f.ensure()

	if err := chkDatPresent(dat, f.inputs); err != nil {
		return OfErr(), err
	}

	ret0, err := f.exe(ctx, dat)
	if err != nil {
		return OfErr(), err
	}

	ret1, err := f.replicate(ret0)
	if err != nil {
		return OfErr(), err
	}

	ret2, err := f.swap(ret1)
	if err != nil {
		return OfErr(), err
	}

	ret3, err := f.select_(ret2)
	if err != nil {
		return OfErr(), err
	}

	ret4, err := f.scope(ret3)
	if err != nil {
		return OfErr(), err
	}

	if err := chkDatPresent(ret4, f.outputs); err != nil {
		return OfErr(), err
	}

	return ret4, nil
}
