package hippo

import (
	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/internal/dot1"
)

func Static(
	dat giraffe.Datum,
) *Fn {
	return M(TryFnOf(func(
		Context,
		giraffe.Datum,
	) (giraffe.Datum, error) {
		return dat, nil
	}))
}

func StaticOf(
	pairs ...giraffe.Tuple,
) (*Fn, error) {
	dat, err := OfN(pairs...)
	if err != nil {
		return nil, err
	}

	return Static(dat), nil
}
