package hippo

import (
	"github.com/hkoosha/giraffe"
)

func Static(
	dat giraffe.Datum,
) *Fn {
	return FnOf(func(
		Context,
		giraffe.Datum,
	) (giraffe.Datum, error) {
		return dat, nil
	})
}

func StaticOf(
	pairs ...giraffe.Tuple,
) (*Fn, error) {
	dat, err := giraffe.OfN(pairs...)
	if err != nil {
		return nil, err
	}

	return Static(dat), nil
}

// =====================================

func SelectRand() *Fn {
	return FnOf(selectRand)
}

func selectRand(
	ctx Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	l, err := dat.Len()
	if err != nil {
		return giraffe.OfErr(), err
	}

	// TODO how to sanely select obj keys? Probably never.
	// TODO chk int range.

	i := ctx.Rand().StdV2().UintN(uint(l)) //nolint:gosec // TODO fix later
	return dat.At(int(i))                  //nolint:gosec // TODO fix later
}

// =====================================

func SelectKey1(
	key giraffe.Query,
) *Fn {
	return FnOf(selectKey{key: key}.exe)
}

type selectKey struct {
	key giraffe.Query
}

func (k selectKey) exe(
	_ Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	key, err := dat.QStr(k.key)
	if err != nil {
		return giraffe.OfErr(), err
	}

	return dat.Get(giraffe.Q(key))
}
