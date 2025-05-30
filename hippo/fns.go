package hippo

import (
	"context"

	"github.com/hkoosha/giraffe"
)

func Static(
	dat giraffe.Datum,
) *Fn {
	return Of(func(
		context.Context,
		giraffe.Datum,
	) (giraffe.Datum, error) {
		return dat, nil
	})
}
