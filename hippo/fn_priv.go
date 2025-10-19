package hippo

import (
	"errors"
	"maps"
	"slices"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/t11y"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
	. "github.com/hkoosha/giraffe/internal/dot1"
)

var errInvalidFn = errors.New("invalid fn")

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

// =====================================

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

	return Of(selected), nil
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

	return Of(ret), nil
}

// =====================================

func (f *Fn) clone() *Fn {
	f.ensure()

	if f == nil {
		return nil
	}

	return &Fn{
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

func (f *Fn) call(
	ctx Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	t11y.NonNil(f, f.exe)

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
