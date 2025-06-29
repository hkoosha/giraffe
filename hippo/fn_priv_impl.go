package hippo

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/g11y/gtx"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
	. "github.com/hkoosha/giraffe/internal/dot0"
	. "github.com/hkoosha/giraffe/internal/dot1"
)

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

		for _, i := range into {
			dat, err = dat.Set(i, val)
			if err != nil {
				return OfErr(), err
			}
		}
	}

	return dat, nil
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
		if swapTo, ok := f.swapped[k]; ok {
			k = swapTo
		}

		ret[k] = v
	}

	return Of0(ret), nil
}

func (f *Fn_) scopeOut(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	if f.scopedOut == "" {
		return dat, nil
	}

	return giraffe.Of1(f.scopedOut, dat), nil
}

func (f *Fn_) scopeIn(
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	if f.scopedIn == "" {
		return dat, nil
	}

	if !dat.Has(f.scopedIn) {
		return OfEmpty(), hippoerr.NewFnMissingKeysError(f.scopedIn)
	}

	return dat.Get(f.scopedIn)
}

// =====================================.

func (f *Fn_) call(
	ctx gtx.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	g11y.NonNil(f, f.exe)

	f.ensure()

	dat0, err := f.scopeIn(dat)
	if err != nil {
		return OfErr(), err
	}

	if err0 := chkDatPresent(dat0, f.inputs); err0 != nil {
		return OfErr(), err0
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

	ret4, err := f.scopeOut(ret3)
	if err != nil {
		return OfErr(), err
	}

	if err := chkDatPresent(ret4, f.outputs); err != nil {
		return OfErr(), err
	}

	if f.noOverwriting && dat.HasAll(f.outputs...) {
		return OfEmpty(), nil
	}

	return ret4, nil
}
