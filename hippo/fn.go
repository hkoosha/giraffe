package hippo

import (
	"context"
	"maps"
	"slices"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
	. "github.com/hkoosha/giraffe/internal/dot0"
	. "github.com/hkoosha/giraffe/internal/dot1"
)

type Exe func(
	context.Context,
	giraffe.Datum,
) (giraffe.Datum, error)

// =============================================================================.

func Of(
	exe Exe,
) *Fn {
	return &Fn{
		exe:        exe,
		scoped:     "",
		inputs:     nil,
		outputs:    nil,
		replicated: nil,
		selected:   nil,
		swapped:    nil,
	}
}

type Fn struct {
	exe        Exe
	replicated map[giraffe.Query]giraffe.Query
	swapped    map[giraffe.Query]giraffe.Query
	scoped     giraffe.Query
	inputs     []giraffe.Query
	outputs    []giraffe.Query
	selected   []giraffe.Query
}

func (f *Fn) AndReplicate(
	replicated map[giraffe.Query]giraffe.Query,
) *Fn {
	replicated = maps.Clone(replicated)
	maps.Copy(replicated, f.replicated)

	clone := f.clone()
	clone.replicated = replicated
	return clone
}

func (f *Fn) WithReplicated(
	replicated map[giraffe.Query]giraffe.Query,
) *Fn {
	clone := f.clone()
	clone.replicated = maps.Clone(replicated)
	return clone
}

func (f *Fn) AndSwapping(
	swapping map[giraffe.Query]giraffe.Query,
) *Fn {
	swapping = maps.Clone(swapping)
	maps.Copy(swapping, f.swapped)

	clone := f.clone()
	clone.swapped = swapping
	return clone
}

func (f *Fn) WithSwapping(
	swapping map[giraffe.Query]giraffe.Query,
) *Fn {
	clone := f.clone()
	clone.swapped = maps.Clone(swapping)
	return clone
}

func (f *Fn) AndScope(
	scope giraffe.Query,
) *Fn {
	return f.WithScope(f.scoped.Plus(scope))
}

func (f *Fn) WithScope(
	scope giraffe.Query,
) *Fn {
	clone := f.clone()
	clone.scoped = scope
	return clone
}

func (f *Fn) AndInputs(
	inputs ...giraffe.Query,
) *Fn {
	clone := f.clone()
	clone.inputs = append(slices.Clone(clone.inputs), inputs...)
	return clone
}

func (f *Fn) WithInput(
	inputs ...giraffe.Query,
) *Fn {
	clone := f.clone()
	clone.inputs = slices.Clone(inputs)
	return clone
}

func (f *Fn) AndOutputs(
	outputs ...giraffe.Query,
) *Fn {
	clone := f.clone()
	clone.outputs = append(slices.Clone(clone.outputs), outputs...)
	return clone
}

func (f *Fn) WithOutput(
	outputs ...giraffe.Query,
) *Fn {
	clone := f.clone()
	clone.outputs = slices.Clone(outputs)
	return clone
}

func (f *Fn) Select(
	select_ ...giraffe.Query,
) *Fn {
	clone := f.clone()
	clone.selected = slices.Clone(select_)
	return clone
}

func (f *Fn) SelectAll() *Fn {
	clone := f.clone()
	clone.selected = nil
	return clone
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
		return hippoerr.NewFnMissingKeysError(dat, missing)
	}

	return nil
}

// =====================================.

func (f *Fn) replicate(
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

func (f *Fn) scope(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	if f.scoped == "" {
		return dat, nil
	}

	return giraffe.Of1(f.scoped, dat), nil
}

func (f *Fn) select_(
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

func (f *Fn) swap(
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

func (f *Fn) clone() *Fn {
	if f == nil {
		return nil
	}

	return &Fn{
		exe:        f.exe,
		scoped:     f.scoped,
		inputs:     slices.Clone(f.inputs),
		outputs:    slices.Clone(f.outputs),
		replicated: maps.Clone(f.replicated),
		swapped:    maps.Clone(f.swapped),
		selected:   slices.Clone(f.selected),
	}
}

func (f *Fn) call(
	ctx context.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	g11y.NonNil(f, f.exe)

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
