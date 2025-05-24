package hippo

import (
	"errors"
	"maps"
	"slices"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/t11y"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
)

var errInvalidFn = errors.New("invalid fn")

func allExists(
	dat giraffe.Datum,
	keys []giraffe.Query,
) bool {
	for _, k := range keys {
		if ok, err := dat.Has(k); err != nil || !ok {
			return false
		}
	}

	return true
}

func cloneMapOfSlices[K comparable, S ~[]E, E any](
	gather map[K]S,
) map[K]S {
	cloned := make(map[K]S, len(gather))
	for k, v := range gather {
		cloned[k] = slices.Clone(v)
	}

	return cloned
}

// =====================================

func (f *Fn) copied(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	if len(f.copy) == 0 {
		return dat, nil
	}

	for from, into := range f.copy {
		if ok, err := dat.Has(from); err != nil {
			return dErr, err
		} else if !ok {
			continue
		}

		val, err := dat.Get(from)
		if err != nil {
			return dErr, err
		}

		dat, err = dat.Set(into.WithMake(), val)
		if err != nil {
			return dErr, err
		}
	}

	return dat, nil
}

func (f *Fn) scope(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	if f.scoped == nil {
		return dat, nil
	}

	return giraffe.Of1(*f.scoped, dat), nil
}

func (f *Fn) select_(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	if len(f.selected) == 0 {
		return dat, nil
	}

	selected := make(map[giraffe.Query]giraffe.Datum, len(f.selected))
	for _, k := range f.selected {
		has, err := dat.Has(k)
		if err != nil {
			return dErr, err
		} else if !has {
			continue
		}

		v, err := dat.Get(k)
		if err != nil {
			return dErr, err
		}

		selected[k] = v
	}

	return giraffe.Of(selected), nil
}

/*func (f *Fn) swap(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	if len(f.swapped) == 0 {
		return dat, nil
	}

	if !dat.Type().IsObj() {
		return dErr, EF("expecting an object, got: %s", dat.Type())
	}

	iter, err := dat.Iter2()
	if err != nil {
		return dErr, err
	}

	ret := make(map[giraffe.Query]giraffe.Datum, M(dat.Len()))
	for k, v := range iter {
		k := giraffe.Q(k)
		if swapTo, ok := f.copied[k]; ok {
			k = swapTo
		}

		ret[k] = v
	}

	return giraffe.Of(ret), nil
}*/

// =====================================

func (f *Fn) clone() *Fn {
	f.ensure()

	if f == nil {
		return nil
	}

	return &Fn{
		exe:          f.exe,
		scoped:       f.scoped,
		combine:      cloneMapOfSlices(f.combine),
		inputs:       slices.Clone(f.inputs),
		optionals:    slices.Clone(f.optionals),
		outputs:      slices.Clone(f.outputs),
		copy:         maps.Clone(f.copy),
		selected:     slices.Clone(f.selected),
		skipOnExists: f.skipOnExists,
		skipped:      f.skipped,
		skipWith:     f.skipWith,
		typ:          f.typ.Clone(),
		name:         f.name,

		// args:      slices.Clone(f.args),
		// swapped:      maps.Clone(f.swapped),
	}
}

func (f *Fn) call(
	ctx gtx.Context,
	call Call,
) (giraffe.Datum, error) {
	t11y.NonNil(f, f.exe)
	if f.skipWith != nil {
		return *f.skipWith, nil
	}
	if f.skipped {
		return giraffe.OfEmpty(), nil
	}

	dat := call.Data()

	f.ensure()

	if err := call.CheckPresent(dat, f.inputs); err != nil {
		return dErr, err
	}

	// if err := call.CheckPresent(call.Args(), f.args); err != nil {
	// 	return dErr, err
	// }

	if f.skipOnExists && allExists(dat, f.outputs) {
		return dat, nil
	}

	if len(f.combine) > 0 {
		for into, froms := range f.combine {
			gathered := giraffe.OfEmpty()
			for i, from := range froms {
				d, err := dat.Get(from)
				if err != nil {
					return dErr, err
				}

				switch {
				case i == 0:
					gathered, err = gathered.Set(into, d)
				default:
					gathered, err = gathered.Merge(d)
				}

				if err != nil {
					return dErr, err
				}
			}

			var err error
			dat, err = dat.Merge(gathered)
			if err != nil {
				return dErr, err
			}
		}
	}

	ret0, err := f.exe(ctx, call.WithData(dat))
	if err != nil {
		return dErr, err
	}

	ret1, err := f.copied(ret0)
	if err != nil {
		return dErr, err
	}

	if cErr := call.CheckPresent(ret1, f.outputs); cErr != nil {
		return dErr, cErr
	}

	ret2, err := f.select_(ret1)
	if err != nil {
		return dErr, err
	}

	ret3, err := f.scope(ret2)
	if err != nil {
		return dErr, err
	}

	return ret3, nil
}
