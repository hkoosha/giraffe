package hippo

import (
	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/t11y/dot"
)

func Static(
	dat giraffe.Datum,
) *Fn {
	return M(FnOf(func(
		Context,
		giraffe.Datum,
	) (giraffe.Datum, error) {
		return dat, nil
	}))
}
