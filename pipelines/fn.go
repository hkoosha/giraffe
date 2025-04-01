package pipelines

import (
	"context"

	"github.com/hkoosha/giraffe"
)

type Fn func(
	context.Context,
	giraffe.Datum,
) (giraffe.Datum, error)
