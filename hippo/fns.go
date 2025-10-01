package hippo

import (
	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

func Static(
	dat giraffe.Datum,
) *Fn_ {
	return M(FnOf(func(
		Context,
		giraffe.Datum,
	) (giraffe.Datum, error) {
		return dat, nil
	}))
}
